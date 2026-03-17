use actix_web::{web, HttpResponse};
use serde::Deserialize;

use crate::db::Database;
use crate::error::AppError;
use crate::models::*;
use crate::redis::RedisClient;

#[derive(Debug, Deserialize)]
pub struct InteractionStatus {
    pub is_liked: bool,
    pub is_favorited: bool,
}

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/interactions")
            .route("/like/{video_id}", web::post().to(like_video))
            .route("/unlike/{video_id}", web::post().to(unlike_video))
            .route("/favorite/{video_id}", web::post().to(favorite_video))
            .route("/unfavorite/{video_id}", web::post().to(unfavorite_video))
            .route("/status/{video_id}", web::get().to(get_status))
            .route("/liked", web::get().to(get_liked_videos))
            .route("/favorites", web::get().to(get_favorite_videos)),
    );
}

async fn like_video(
    db: web::Data<Database>,
    _redis: web::Data<RedisClient>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let video_id = path.into_inner();
    
    let existing: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM likes WHERE user_id = $1 AND target_id = $2 AND type = 'video'"
    )
    .bind(*user_id)
    .bind(video_id)
    .fetch_one(db.pool())
    .await?;

    if existing > 0 {
        return Err(AppError::Conflict("Already liked".into()));
    }

    sqlx::query("INSERT INTO likes (user_id, target_id, type) VALUES ($1, $2, 'video')")
        .bind(*user_id)
        .bind(video_id)
        .execute(db.pool())
        .await?;

    sqlx::query("UPDATE videos SET like_count = like_count + 1 WHERE id = $1")
        .bind(video_id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({ "success": true })))
}

async fn unlike_video(
    db: web::Data<Database>,
    _redis: web::Data<RedisClient>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let video_id = path.into_inner();
    
    let result = sqlx::query(
        "DELETE FROM likes WHERE user_id = $1 AND target_id = $2 AND type = 'video'"
    )
    .bind(*user_id)
    .bind(video_id)
    .execute(db.pool())
    .await?;

    if result.rows_affected() == 0 {
        return Err(AppError::NotFound("Not liked".into()));
    }

    sqlx::query("UPDATE videos SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1")
        .bind(video_id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({ "success": true })))
}

async fn favorite_video(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let video_id = path.into_inner();
    
    let existing: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM favorites WHERE user_id = $1 AND video_id = $2"
    )
    .bind(*user_id)
    .bind(video_id)
    .fetch_one(db.pool())
    .await?;

    if existing > 0 {
        return Err(AppError::Conflict("Already favorited".into()));
    }

    sqlx::query("INSERT INTO favorites (user_id, video_id) VALUES ($1, $2)")
        .bind(*user_id)
        .bind(video_id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({ "success": true })))
}

async fn unfavorite_video(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let video_id = path.into_inner();
    
    let result = sqlx::query(
        "DELETE FROM favorites WHERE user_id = $1 AND video_id = $2"
    )
    .bind(*user_id)
    .bind(video_id)
    .execute(db.pool())
    .await?;

    if result.rows_affected() == 0 {
        return Err(AppError::NotFound("Not favorited".into()));
    }

    Ok(HttpResponse::Ok().json(serde_json::json!({ "success": true })))
}

async fn get_status(
    db: web::Data<Database>,
    _redis: web::Data<RedisClient>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let video_id = path.into_inner();
    
    let is_liked: bool = sqlx::query_scalar(
        "SELECT COUNT(*) > 0 FROM likes WHERE user_id = $1 AND target_id = $2 AND type = 'video'"
    )
    .bind(*user_id)
    .bind(video_id)
    .fetch_one(db.pool())
    .await?;

    let is_favorited: bool = sqlx::query_scalar(
        "SELECT COUNT(*) > 0 FROM favorites WHERE user_id = $1 AND video_id = $2"
    )
    .bind(*user_id)
    .bind(video_id)
    .fetch_one(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(InteractionStatus { is_liked, is_favorited }))
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

async fn get_liked_videos(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let offset = (query.page - 1) * query.page_size;
    
    let videos: Vec<VideoResponse> = sqlx::query_as(
        r#"
        SELECT v.id, v.user_id, v.title, v.description, v.status::text, v.visibility::text, 
               v.duration_seconds, v.thumbnail_url, v.category_id, v.view_count, v.like_count, v.comment_count, v.published_at, v.created_at
        FROM videos v
        JOIN likes l ON v.id = l.target_id
        WHERE l.user_id = $1 AND l.type = 'video' AND v.status = 'Published'
        ORDER BY l.created_at DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(*user_id)
    .bind(query.page_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(videos))
}

async fn get_favorite_videos(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let offset = (query.page - 1) * query.page_size;
    
    let videos: Vec<VideoResponse> = sqlx::query_as(
        r#"
        SELECT v.id, v.user_id, v.title, v.description, v.status::text, v.visibility::text, 
               v.duration_seconds, v.thumbnail_url, v.category_id, v.view_count, v.like_count, v.comment_count, v.published_at, v.created_at
        FROM videos v
        JOIN favorites f ON v.id = f.video_id
        WHERE f.user_id = $1 AND v.status = 'Published'
        ORDER BY f.created_at DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(*user_id)
    .bind(query.page_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(videos))
}
