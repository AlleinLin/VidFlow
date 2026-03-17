use actix_web::{web, HttpResponse};
use serde::Deserialize;

use crate::db::Database;
use crate::error::AppError;
use crate::models::*;

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/search")
            .route("/", web::get().to(search_all))
            .route("/videos", web::get().to(search_videos))
            .route("/users", web::get().to(search_users))
    );
}

#[derive(Debug, Deserialize)]
struct SearchQuery {
    q: String,
    #[serde(default = "default_page")]
    page: i32,
    #[serde(default = "default_page_size")]
    page_size: i32,
}

fn default_page() -> i32 { 1 }
fn default_page_size() -> i32 { 20 }

async fn search_all(
    db: web::Data<Database>,
    query: web::Query<SearchQuery>,
) -> Result<HttpResponse, AppError> {
    let half_size = query.page_size / 2;
    let offset = (query.page - 1) * half_size;
    let pattern = format!("%{}%", query.q);
    
    let videos: Vec<VideoResponse> = sqlx::query_as(
        r#"
        SELECT id, user_id, title, description, status::text, visibility::text, 
               duration_seconds, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at
        FROM videos 
        WHERE status = 'Published' AND (title ILIKE $1 OR description ILIKE $1)
        ORDER BY view_count DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(&pattern)
    .bind(half_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    #[derive(Debug, Clone, FromRow, Serialize)]
    struct UserSearchResult {
        id: i64,
        username: String,
        display_name: String,
        avatar_url: Option<String>,
        bio: Option<String>,
        follower_count: i64,
    }

    let users: Vec<UserSearchResult> = sqlx::query_as(
        r#"
        SELECT id, username, display_name, avatar_url, bio, follower_count
        FROM users 
        WHERE status = 'Active' AND (username ILIKE $1 OR display_name ILIKE $1)
        ORDER BY follower_count DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(&pattern)
    .bind(half_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "videos": videos,
        "users": users,
        "page": query.page,
        "page_size": query.page_size
    })))
}

async fn search_videos(
    db: web::Data<Database>,
    query: web::Query<SearchQuery>,
) -> Result<HttpResponse, AppError> {
    let offset = (query.page - 1) * query.page_size;
    let pattern = format!("%{}%", query.q);
    
    let total: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM videos WHERE status = 'Published' AND (title ILIKE $1 OR description ILIKE $1)"
    )
    .bind(&pattern)
    .fetch_one(db.pool())
    .await?;

    let videos: Vec<VideoResponse> = sqlx::query_as(
        r#"
        SELECT id, user_id, title, description, status::text, visibility::text, 
               duration_seconds, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at
        FROM videos 
        WHERE status = 'Published' AND (title ILIKE $1 OR description ILIKE $1)
        ORDER BY view_count DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(&pattern)
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

async fn search_users(
    db: web::Data<Database>,
    query: web::Query<SearchQuery>,
) -> Result<HttpResponse, AppError> {
    let offset = (query.page - 1) * query.page_size;
    let pattern = format!("%{}%", query.q);
    
    #[derive(Debug, Clone, FromRow, Serialize)]
    struct UserSearchResult {
        id: i64,
        username: String,
        display_name: String,
        avatar_url: Option<String>,
        bio: Option<String>,
        follower_count: i64,
    }

    let total: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM users WHERE status = 'Active' AND (username ILIKE $1 OR display_name ILIKE $1)"
    )
    .bind(&pattern)
    .fetch_one(db.pool())
    .await?;

    let users: Vec<UserSearchResult> = sqlx::query_as(
        r#"
        SELECT id, username, display_name, avatar_url, bio, follower_count
        FROM users 
        WHERE status = 'Active' AND (username ILIKE $1 OR display_name ILIKE $1)
        ORDER BY follower_count DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(&pattern)
    .bind(query.page_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "results": users,
        "total": total,
        "page": query.page,
        "page_size": query.page_size
    })))
}
