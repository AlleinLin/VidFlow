use actix_web::{web, HttpResponse};
use serde::Deserialize;
use argon2::{password_hash::SaltString, Argon2, PasswordHash, PasswordHasher, PasswordVerifier};
use rand::rngs::OsRng;
use chrono::{Duration, Utc};
use jsonwebtoken::{decode, encode, DecodingKey, EncodingKey, Header, Validation};

use crate::config::AppConfig;
use crate::db::Database;
use crate::error::AppError;
use crate::models::*;
use crate::redis::RedisClient;
use crate::services::{Claims, TokenPair};

#[derive(Debug, Deserialize)]
pub struct RegisterRequest {
    pub username: String,
    pub email: String,
    pub password: String,
    pub display_name: String,
}

#[derive(Debug, Deserialize)]
pub struct LoginRequest {
    pub email: String,
    pub password: String,
}

#[derive(Debug, serde::Serialize)]
pub struct LoginResponse {
    pub token: TokenPair,
    pub user: UserProfile,
}

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/auth")
            .route("/register", web::post().to(register))
            .route("/login", web::post().to(login))
            .route("/refresh", web::post().to(refresh_token))
            .route("/logout", web::post().to(logout)),
    );
}

async fn register(
    db: web::Data<Database>,
    _redis: web::Data<RedisClient>,
    config: web::Data<AppConfig>,
    body: web::Json<RegisterRequest>,
) -> Result<HttpResponse, AppError> {
    let existing: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM users WHERE email = $1 OR username = $2",
    )
    .bind(&body.email)
    .bind(&body.username)
    .fetch_one(db.pool())
    .await?;

    if existing > 0 {
        return Err(AppError::Conflict("Email or username already exists".into()));
    }

    let password_hash = hash_password(&body.password)?;

    let user = sqlx::query_as::<_, User>(
        r#"
        INSERT INTO users (username, email, password_hash, display_name, role, status)
        VALUES ($1, $2, $3, $4, 'User', 'Active')
        RETURNING *
        "#,
    )
    .bind(&body.username)
    .bind(&body.email)
    .bind(&password_hash)
    .bind(&body.display_name)
    .fetch_one(db.pool())
    .await?;

    let tokens = generate_token_pair(&user, &config.jwt)?;

    Ok(HttpResponse::Ok().json(LoginResponse {
        token: tokens,
        user: UserProfile::from(user),
    }))
}

async fn login(
    db: web::Data<Database>,
    _redis: web::Data<RedisClient>,
    config: web::Data<AppConfig>,
    body: web::Json<LoginRequest>,
) -> Result<HttpResponse, AppError> {
    let user = sqlx::query_as::<_, User>("SELECT * FROM users WHERE email = $1")
        .bind(&body.email)
        .fetch_optional(db.pool())
        .await?
        .ok_or_else(|| AppError::Unauthorized("Invalid email or password".into()))?;

    if !matches!(user.status, UserStatus::Active) {
        return Err(AppError::Unauthorized("Account is not active".into()));
    }

    verify_password(&body.password, &user.password_hash)?;

    let tokens = generate_token_pair(&user, &config.jwt)?;

    sqlx::query("UPDATE users SET last_login_at = NOW() WHERE id = $1")
        .bind(user.id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(LoginResponse {
        token: tokens,
        user: UserProfile::from(user),
    }))
}

#[derive(Debug, Deserialize)]
struct RefreshRequest {
    refresh_token: String,
}

async fn refresh_token(
    db: web::Data<Database>,
    config: web::Data<AppConfig>,
    body: web::Json<RefreshRequest>,
) -> Result<HttpResponse, AppError> {
    let claims = validate_token(&body.refresh_token, &config.jwt.secret)?;

    let user = sqlx::query_as::<_, User>("SELECT * FROM users WHERE id = $1")
        .bind(claims.user_id)
        .fetch_optional(db.pool())
        .await?
        .ok_or_else(|| AppError::NotFound("User not found".into()))?;

    if !matches!(user.status, UserStatus::Active) {
        return Err(AppError::Unauthorized("Account is not active".into()));
    }

    let tokens = generate_token_pair(&user, &config.jwt)?;

    Ok(HttpResponse::Ok().json(tokens))
}

async fn logout(
    _db: web::Data<Database>,
    _user_id: web::ReqData<i64>,
) -> Result<HttpResponse, AppError> {
    Ok(HttpResponse::Ok().json(serde_json::json!({ "success": true })))
}

fn hash_password(password: &str) -> Result<String, AppError> {
    let salt = SaltString::generate(&mut OsRng);
    let argon2 = Argon2::default();
    argon2
        .hash_password(password.as_bytes(), &salt)
        .map(|hash| hash.to_string())
        .map_err(|_| AppError::Internal("Failed to hash password".into()))
}

fn verify_password(password: &str, hash: &str) -> Result<(), AppError> {
    let parsed_hash = PasswordHash::new(hash)
        .map_err(|_| AppError::Internal("Invalid password hash".into()))?;
    Argon2::default()
        .verify_password(password.as_bytes(), &parsed_hash)
        .map_err(|_| AppError::Unauthorized("Invalid email or password".into()))
}

fn generate_token_pair(user: &User, config: &crate::config::JwtConfig) -> Result<TokenPair, AppError> {
    let now = Utc::now();
    let encoding_key = EncodingKey::from_secret(config.secret.as_bytes());

    let access_claims = Claims {
        sub: user.id.to_string(),
        user_id: user.id,
        username: user.username.clone(),
        role: match user.role {
            UserRole::User => "USER".to_string(),
            UserRole::Creator => "CREATOR".to_string(),
            UserRole::Moderator => "MODERATOR".to_string(),
            UserRole::Admin => "ADMIN".to_string(),
        },
        exp: (now + Duration::seconds(config.access_token_ttl)).timestamp(),
        iat: now.timestamp(),
        iss: config.issuer.clone(),
    };

    let access_token = encode(&Header::default(), &access_claims, &encoding_key)
        .map_err(|_| AppError::Internal("Failed to generate access token".into()))?;

    let refresh_claims = Claims {
        sub: user.id.to_string(),
        user_id: user.id,
        username: user.username.clone(),
        role: access_claims.role.clone(),
        exp: (now + Duration::seconds(config.refresh_token_ttl)).timestamp(),
        iat: now.timestamp(),
        iss: config.issuer.clone(),
    };

    let refresh_token = encode(&Header::default(), &refresh_claims, &encoding_key)
        .map_err(|_| AppError::Internal("Failed to generate refresh token".into()))?;

    Ok(TokenPair {
        access_token,
        refresh_token,
        expires_in: config.access_token_ttl,
        token_type: "Bearer".to_string(),
    })
}

fn validate_token(token: &str, secret: &str) -> Result<Claims, AppError> {
    let decoding_key = DecodingKey::from_secret(secret.as_bytes());
    let token_data = decode::<Claims>(token, &decoding_key, &Validation::default())
        .map_err(|_| AppError::Unauthorized("Invalid token".into()))?;

    Ok(token_data.claims)
}
