use crate::config::JwtConfig;
use crate::error::AppError;
use crate::models::*;
use crate::{db::Database, redis::RedisClient};
use argon2::{password_hash::SaltString, Argon2, PasswordHash, PasswordHasher, PasswordVerifier};
use chrono::{Duration, Utc};
use jsonwebtoken::{decode, encode, DecodingKey, EncodingKey, Header, Validation};
use rand::rngs::OsRng;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    pub sub: String,
    pub user_id: i64,
    pub username: String,
    pub role: String,
    pub exp: i64,
    pub iat: i64,
    pub iss: String,
}

pub struct AuthService {
    db: Database,
    redis: RedisClient,
    jwt_config: JwtConfig,
}

impl AuthService {
    pub fn new(db: Database, redis: RedisClient, jwt_config: JwtConfig) -> Self {
        Self {
            db,
            redis,
            jwt_config,
        }
    }

    pub async fn register(
        &self,
        username: &str,
        email: &str,
        password: &str,
        display_name: &str,
    ) -> Result<User, AppError> {
        let existing = sqlx::query_scalar::<_, i64>(
            "SELECT COUNT(*) FROM users WHERE email = $1 OR username = $2",
        )
        .bind(email)
        .bind(username)
        .fetch_one(self.db.pool())
        .await?;

        if existing > 0 {
            return Err(AppError::Conflict("Email or username already exists".into()));
        }

        let password_hash = hash_password(password)?;

        let user = sqlx::query_as::<_, User>(
            r#"
            INSERT INTO users (username, email, password_hash, display_name, role, status)
            VALUES ($1, $2, $3, $4, 'User', 'Active')
            RETURNING *
            "#,
        )
        .bind(username)
        .bind(email)
        .bind(&password_hash)
        .bind(display_name)
        .fetch_one(self.db.pool())
        .await?;

        Ok(user)
    }

    pub async fn login(&self, email: &str, password: &str) -> Result<(User, TokenPair), AppError> {
        let user = sqlx::query_as::<_, User>("SELECT * FROM users WHERE email = $1")
            .bind(email)
            .fetch_optional(self.db.pool())
            .await?
            .ok_or_else(|| AppError::Unauthorized("Invalid email or password".into()))?;

        if !matches!(user.status, UserStatus::Active) {
            return Err(AppError::Unauthorized("Account is not active".into()));
        }

        verify_password(password, &user.password_hash)?;

        let tokens = generate_token_pair(&user, &self.jwt_config)?;

        sqlx::query("UPDATE users SET last_login_at = NOW() WHERE id = $1")
            .bind(user.id)
            .execute(self.db.pool())
            .await?;

        Ok((user, tokens))
    }

    pub async fn validate_token(&self, token: &str) -> Result<Claims, AppError> {
        let decoding_key = DecodingKey::from_secret(self.jwt_config.secret.as_bytes());
        let token_data = decode::<Claims>(token, &decoding_key, &Validation::default())
            .map_err(|_| AppError::Unauthorized("Invalid token".into()))?;

        Ok(token_data.claims)
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct TokenPair {
    pub access_token: String,
    pub refresh_token: String,
    pub expires_in: i64,
    pub token_type: String,
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

fn generate_token_pair(user: &User, config: &JwtConfig) -> Result<TokenPair, AppError> {
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
        role: match user.role {
            UserRole::User => "USER".to_string(),
            UserRole::Creator => "CREATOR".to_string(),
            UserRole::Moderator => "MODERATOR".to_string(),
            UserRole::Admin => "ADMIN".to_string(),
        },
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
