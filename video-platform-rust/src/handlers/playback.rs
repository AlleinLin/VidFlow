use actix_web::{web, HttpResponse};
use serde::Deserialize;
use chrono::Utc;

use crate::db::Database;
use crate::error::AppError;
use crate::models::*;

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/playback")
            .route("/progress", web::post().to(update_progress))
            .route("/progress/{video_id}", web::get().to(get_progress))
            .route("/history", web::get().to(get_watch_history))
            .route("/continue-watching", web::get().to(get_continue_watching))
            .route("/history/{video_id}", web::delete().to(delete_watch_history))
            .route("/history", web::delete().to(clear_watch_history))
    );
}

async fn update_progress(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    req: web::Json<UpdateProgressRequest>,
) -> Result<HttpResponse, AppError> {
    let progress = if req.duration > 0.0 { req.position / req.duration } else { 0.0 };
    let completed = progress >= 0.95;
    let now = Utc::now();
    
    let existing: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM watch_history WHERE user_id = $1 AND video_id = $2"
    )
    .bind(*user_id)
    .bind(req.video_id)
    .fetch_one(db.pool())
    .await?;

    if existing > 0 {
        sqlx::query(
            r#"
            UPDATE watch_history 
            SET watch_duration = watch_duration + $3,
                watch_progress = $4,
                last_position = $5,
                completed = $6,
                watched_at = $7,
                updated_at = $7
            WHERE user_id = $1 AND video_id = $2
            "#,
        )
        .bind(*user_id)
        .bind(req.video_id)
        .bind(req.watch_duration)
        .bind(progress)
        .bind(req.position)
        .bind(completed)
        .bind(now)
        .execute(db.pool())
        .await?;
    } else {
        sqlx::query(
            r#"
            INSERT INTO watch_history (user_id, video_id, watch_duration, watch_progress, last_position, completed, watched_at, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $7, $7)
            "#,
        )
        .bind(*user_id)
        .bind(req.video_id)
        .bind(req.watch_duration)
        .bind(progress)
        .bind(req.position)
        .bind(completed)
        .bind(now)
        .execute(db.pool())
        .await?;

        sqlx::query("UPDATE videos SET view_count = view_count + 1 WHERE id = $1")
            .bind(req.video_id)
            .execute(db.pool())
            .await?;
    }

    Ok(HttpResponse::Ok().json(serde_json::json!({ "message": "Progress updated successfully" })))
}

async fn get_progress(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let video_id = path.into_inner();
    
    let history: Option<WatchHistory> = sqlx::query_as(
        "SELECT id, user_id, video_id, watch_duration, watch_progress, last_position, completed, watched_at, created_at, updated_at FROM watch_history WHERE user_id = $1 AND video_id = $2"
    )
    .bind(*user_id)
    .bind(video_id)
    .fetch_optional(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(history))
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

async fn get_watch_history(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let offset = (query.page - 1) * query.page_size;
    
    let histories: Vec<WatchHistoryResponse> = sqlx::query_as(
        r#"
        SELECT wh.id, wh.video_id, v.title as video_title, v.thumbnail_url as video_thumbnail,
               wh.watch_progress, wh.last_position, wh.completed, wh.watched_at
        FROM watch_history wh
        JOIN videos v ON wh.video_id = v.id
        WHERE wh.user_id = $1
        ORDER BY wh.watched_at DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(*user_id)
    .bind(query.page_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    let total: i64 = sqlx::query_scalar("SELECT COUNT(*) FROM watch_history WHERE user_id = $1")
        .bind(*user_id)
        .fetch_one(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "histories": histories,
        "total": total,
        "page": query.page,
        "page_size": query.page_size
    })))
}

#[derive(Debug, Deserialize)]
struct LimitQuery {
    #[serde(default = "default_limit")]
    limit: i32,
}

fn default_limit() -> i32 { 10 }

async fn get_continue_watching(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    query: web::Query<LimitQuery>,
) -> Result<HttpResponse, AppError> {
    let histories: Vec<WatchHistoryResponse> = sqlx::query_as(
        r#"
        SELECT wh.id, wh.video_id, v.title as video_title, v.thumbnail_url as video_thumbnail,
               wh.watch_progress, wh.last_position, wh.completed, wh.watched_at
        FROM watch_history wh
        JOIN videos v ON wh.video_id = v.id
        WHERE wh.user_id = $1 AND wh.completed = false AND wh.watch_progress > 0.05
        ORDER BY wh.watched_at DESC
        LIMIT $2
        "#,
    )
    .bind(*user_id)
    .bind(query.limit)
    .fetch_all(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(histories))
}

async fn delete_watch_history(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let video_id = path.into_inner();
    
    sqlx::query("DELETE FROM watch_history WHERE user_id = $1 AND video_id = $2")
        .bind(*user_id)
        .bind(video_id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::NoContent().finish())
}

async fn clear_watch_history(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
) -> Result<HttpResponse, AppError> {
    sqlx::query("DELETE FROM watch_history WHERE user_id = $1")
        .bind(*user_id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::NoContent().finish())
}
