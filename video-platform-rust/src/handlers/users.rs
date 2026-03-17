use actix_web::{web, HttpResponse};
use serde::Deserialize;

use crate::db::Database;
use crate::error::AppError;
use crate::models::*;
use crate::redis::RedisClient;

#[derive(Debug, Deserialize)]
pub struct UpdateProfileRequest {
    pub display_name: Option<String>,
    pub bio: Option<String>,
    pub avatar_url: Option<String>,
}

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/users")
            .route("/{id}", web::get().to(get_profile))
            .route("/{id}", web::put().to(update_profile))
            .route("/{id}/follow", web::post().to(follow))
            .route("/{id}/unfollow", web::post().to(unfollow))
            .route("/{id}/followers", web::get().to(get_followers))
            .route("/{id}/following", web::get().to(get_following))
            .route("/{id}/videos", web::get().to(get_user_videos)),
    );
}

async fn get_profile(
    db: web::Data<Database>,
    redis: web::Data<RedisClient>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let user_id = path.into_inner();
    let cache_key = crate::redis::keys::user_profile(user_id);
    
    if let Some(cached) = redis.get::<UserProfile>(&cache_key).await? {
        return Ok(HttpResponse::Ok().json(cached));
    }

    let user = sqlx::query_as::<_, User>("SELECT * FROM users WHERE id = $1")
        .bind(user_id)
        .fetch_optional(db.pool())
        .await?
        .ok_or_else(|| AppError::NotFound("User not found".into()))?;

    let profile = UserProfile::from(user);
    let _ = redis.set(&cache_key, &profile, std::time::Duration::from_secs(600)).await;

    Ok(HttpResponse::Ok().json(profile))
}

async fn update_profile(
    db: web::Data<Database>,
    redis: web::Data<RedisClient>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
    body: web::Json<UpdateProfileRequest>,
) -> Result<HttpResponse, AppError> {
    let target_id = path.into_inner();
    
    if *user_id != target_id {
        return Err(AppError::Unauthorized("Cannot update other user's profile".into()));
    }

    if body.display_name.is_some() {
        sqlx::query("UPDATE users SET display_name = $1 WHERE id = $2")
            .bind(&body.display_name)
            .bind(target_id)
            .execute(db.pool())
            .await?;
    }

    if body.bio.is_some() {
        sqlx::query("UPDATE users SET bio = $1 WHERE id = $2")
            .bind(&body.bio)
            .bind(target_id)
            .execute(db.pool())
            .await?;
    }

    if body.avatar_url.is_some() {
        sqlx::query("UPDATE users SET avatar_url = $1 WHERE id = $2")
            .bind(&body.avatar_url)
            .bind(target_id)
            .execute(db.pool())
            .await?;
    }

    let _ = redis.del(&crate::redis::keys::user_profile(target_id)).await;

    let user = sqlx::query_as::<_, User>("SELECT * FROM users WHERE id = $1")
        .bind(target_id)
        .fetch_one(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(UserProfile::from(user)))
}

async fn follow(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let following_id = path.into_inner();
    
    if *user_id == following_id {
        return Err(AppError::BadRequest("Cannot follow yourself".into()));
    }

    let existing: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM user_follows WHERE follower_id = $1 AND following_id = $2"
    )
    .bind(*user_id)
    .bind(following_id)
    .fetch_one(db.pool())
    .await?;

    if existing > 0 {
        return Err(AppError::Conflict("Already following".into()));
    }

    sqlx::query("INSERT INTO user_follows (follower_id, following_id) VALUES ($1, $2)")
        .bind(*user_id)
        .bind(following_id)
        .execute(db.pool())
        .await?;

    sqlx::query("UPDATE users SET following_count = following_count + 1 WHERE id = $1")
        .bind(*user_id)
        .execute(db.pool())
        .await?;

    sqlx::query("UPDATE users SET follower_count = follower_count + 1 WHERE id = $1")
        .bind(following_id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({ "success": true })))
}

async fn unfollow(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let following_id = path.into_inner();

    let result = sqlx::query("DELETE FROM user_follows WHERE follower_id = $1 AND following_id = $2")
        .bind(*user_id)
        .bind(following_id)
        .execute(db.pool())
        .await?;

    if result.rows_affected() == 0 {
        return Err(AppError::NotFound("Not following".into()));
    }

    sqlx::query("UPDATE users SET following_count = GREATEST(following_count - 1, 0) WHERE id = $1")
        .bind(*user_id)
        .execute(db.pool())
        .await?;

    sqlx::query("UPDATE users SET follower_count = GREATEST(follower_count - 1, 0) WHERE id = $1")
        .bind(following_id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({ "success": true })))
}

#[derive(Debug, Deserialize)]
struct PaginationQuery {
    #[serde(default = "default_page")]
    page: i32,
    #[serde(default = "default_page_size")]
    page_size: i32,
}

fn default_page() -> i32 { 1 }
fn default_page_size() -> i32 { 20 }

async fn get_followers(
    db: web::Data<Database>,
    path: web::Path<i64>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let user_id = path.into_inner();
    let offset = (query.page - 1) * query.page_size;

    let users: Vec<UserProfile> = sqlx::query_as(
        r#"
        SELECT u.id, u.username, u.email, u.display_name, u.avatar_url, u.bio, 
               u.role::text, u.status::text, u.follower_count, u.following_count, u.created_at
        FROM users u
        JOIN user_follows uf ON u.id = uf.follower_id
        WHERE uf.following_id = $1
        ORDER BY uf.created_at DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(user_id)
    .bind(query.page_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(users))
}

async fn get_following(
    db: web::Data<Database>,
    path: web::Path<i64>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let user_id = path.into_inner();
    let offset = (query.page - 1) * query.page_size;

    let users: Vec<UserProfile> = sqlx::query_as(
        r#"
        SELECT u.id, u.username, u.email, u.display_name, u.avatar_url, u.bio, 
               u.role::text, u.status::text, u.follower_count, u.following_count, u.created_at
        FROM users u
        JOIN user_follows uf ON u.id = uf.following_id
        WHERE uf.follower_id = $1
        ORDER BY uf.created_at DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(user_id)
    .bind(query.page_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(users))
}

async fn get_user_videos(
    db: web::Data<Database>,
    path: web::Path<i64>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let user_id = path.into_inner();
    let offset = (query.page - 1) * query.page_size;

    let total: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM videos WHERE user_id = $1 AND status = 'Published'"
    )
    .bind(user_id)
    .fetch_one(db.pool())
    .await?;

    let videos: Vec<VideoResponse> = sqlx::query_as(
        r#"
        SELECT id, user_id, title, description, status::text, visibility::text, 
               duration_seconds, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at
        FROM videos 
        WHERE user_id = $1 AND status = 'Published'
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(user_id)
    .bind(query.page_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(VideoListResponse {
        videos,
        total,
        page: query.page,
        page_size: query.page_size,
    }))
}
