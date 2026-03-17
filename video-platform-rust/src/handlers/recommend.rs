use actix_web::{web, HttpResponse};
use serde::Deserialize;

use crate::db::Database;
use crate::error::AppError;
use crate::models::*;
use crate::redis::RedisClient;

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/recommend")
            .route("/home", web::get().to(get_home_feed))
            .route("/related/{video_id}", web::get().to(get_related))
            .route("/search", web::get().to(search)),
    );
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

async fn get_home_feed(
    db: web::Data<Database>,
    redis: web::Data<RedisClient>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let cache_key = crate::redis::keys::recommend_home(0);
    
    if let Some(cached_ids) = redis.get::<Vec<i64>>(&cache_key).await? {
        let videos = fetch_videos_by_ids(&db, &cached_ids, query.page_size).await?;
        return Ok(HttpResponse::Ok().json(VideoListResponse {
            videos,
            total: cached_ids.len() as i64,
            page: query.page,
            page_size: query.page_size,
        }));
    }

    let videos: Vec<VideoResponse> = sqlx::query_as(
        r#"
        SELECT id, user_id, title, description, status::text, visibility::text, 
               duration_seconds, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at
        FROM videos 
        WHERE status = 'Published' 
        ORDER BY (view_count * 0.3 + like_count * 0.7) DESC
        LIMIT $1
        "#,
    )
    .bind(query.page_size)
    .fetch_all(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(VideoListResponse {
        videos,
        total: videos.len() as i64,
        page: query.page,
        page_size: query.page_size,
    }))
}

async fn get_related(
    db: web::Data<Database>,
    path: web::Path<i64>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let video_id = path.into_inner();
    
    let category_id: Option<i64> = sqlx::query_scalar(
        "SELECT category_id FROM videos WHERE id = $1"
    )
    .bind(video_id)
    .fetch_optional(db.pool())
    .await?;

    let videos: Vec<VideoResponse> = if let Some(cat_id) = category_id {
        sqlx::query_as(
            r#"
            SELECT id, user_id, title, description, status::text, visibility::text, 
                   duration_seconds, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at
            FROM videos 
            WHERE status = 'Published' AND id != $1 AND category_id = $2
            ORDER BY view_count DESC
            LIMIT $3
            "#,
        )
        .bind(video_id)
        .bind(cat_id)
        .bind(query.page_size)
        .fetch_all(db.pool())
        .await?
    } else {
        sqlx::query_as(
            r#"
            SELECT id, user_id, title, description, status::text, visibility::text, 
                   duration_seconds, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at
            FROM videos 
            WHERE status = 'Published' AND id != $1
            ORDER BY view_count DESC
            LIMIT $2
            "#,
        )
        .bind(video_id)
        .bind(query.page_size)
        .fetch_all(db.pool())
        .await?
    };

    Ok(HttpResponse::Ok().json(videos))
}

#[derive(Debug, Deserialize)]
struct SearchQuery {
    q: String,
    #[serde(default = "default_page")]
    page: i32,
    #[serde(default = "default_page_size")]
    page_size: i32,
}

async fn search(
    db: web::Data<Database>,
    query: web::Query<SearchQuery>,
) -> Result<HttpResponse, AppError> {
    let offset = (query.page - 1) * query.page_size;
    let pattern = format!("%{}%", query.q);
    
    let total: i64 = sqlx::query_scalar(
        r#"
        SELECT COUNT(*) FROM videos 
        WHERE status = 'Published' AND (title ILIKE $1 OR description ILIKE $1)
        "#,
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

async fn fetch_videos_by_ids(db: &Database, ids: &[i64], limit: i32) -> Result<Vec<VideoResponse>, AppError> {
    let ids_slice: Vec<i64> = ids.iter().take(limit).copied().collect();
    
    let videos: Vec<VideoResponse> = sqlx::query_as(
        r#"
        SELECT id, user_id, title, description, status::text, visibility::text, 
               duration_seconds, thumbnail_url, category_id, view_count, like_count, comment_count, published_at, created_at
        FROM videos 
        WHERE id = ANY($1) AND status = 'Published'
        "#,
    )
    .bind(&ids_slice)
    .fetch_all(db.pool())
    .await?;

    Ok(videos)
}
