use actix_web::{web, HttpResponse};
use serde::Deserialize;

use crate::db::Database;
use crate::error::AppError;
use crate::models::*;
use crate::redis::RedisClient;

#[derive(Debug, Deserialize)]
pub struct CommentResponse {
    pub id: i64,
    pub video_id: i64,
    pub user_id: i64,
    pub parent_id: Option<i64>,
    pub root_id: Option<i64>,
    pub content: String,
    pub like_count: i32,
    pub status: String,
    pub created_at: String,
}

#[derive(Debug, Deserialize)]
pub struct CommentListResponse {
    pub comments: Vec<CommentResponse>,
    pub total: i64,
}

#[derive(Debug, Deserialize)]
pub struct CreateCommentRequest {
    pub video_id: i64,
    pub content: String,
    pub parent_id: Option<i64>,
}

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/comments")
            .route("/video/{video_id}", web::get().to(get_by_video))
            .route("", web::post().to(create))
            .route("/{id}", web::delete().to(delete))
            .route("/{id}/like", web::post().to(like))
            .route("/{id}/unlike", web::post().to(unlike)),
    );
}

async fn get_by_video(
    db: web::Data<Database>,
    path: web::Path<i64>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let offset = (query.page - 1) * query.page_size;
    
    let comments: Vec<CommentResponse> = sqlx::query_as(
        r#"
        SELECT c.id, c.video_id, c.user_id, c.parent_id, c.root_id, c.content, c.like_count, c.status, c.created_at
        FROM comments c
        WHERE c.video_id = $1 AND c.status = 'visible' AND c.parent_id IS NULL
        ORDER BY c.created_at DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(path.into_inner())
    .bind(query.page_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    let total: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM comments WHERE video_id = $1 AND status = 'visible' AND parent_id IS NULL"
    )
    .bind(path.into_inner())
    .fetch_one(db.pool())
    .await?;

    Ok(HttpResponse::Ok().json(CommentListResponse { comments, total }))
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

async fn create(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    body: web::Json<CreateCommentRequest>,
) -> Result<HttpResponse, AppError> {
    let comment = sqlx::query_as::<_, CommentResponse>(
        r#"
        INSERT INTO comments (video_id, user_id, content, parent_id, status)
        VALUES ($1, $2, $3, $4, 'visible')
        RETURNING id, video_id, user_id, parent_id, root_id, content, like_count, status, created_at
        "#,
    )
    .bind(body.video_id)
    .bind(*user_id)
    .bind(&body.content)
    .bind(body.parent_id)
    .fetch_one(db.pool())
    .await?;

    sqlx::query("UPDATE videos SET comment_count = comment_count + 1 WHERE id = $1")
        .bind(body.video_id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(comment))
}

async fn delete(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let result = sqlx::query(
        "UPDATE comments SET status = 'deleted' WHERE id = $1 AND user_id = $2"
    )
    .bind(path.into_inner())
    .bind(*user_id)
    .execute(db.pool())
    .await?;

    if result.rows_affected() == 0 {
        return Err(AppError::NotFound("Comment not found or unauthorized".into()));
    }

    Ok(HttpResponse::Ok().finish())
}

async fn like(
    db: web::Data<Database>,
    _redis: web::Data<RedisClient>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let existing: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM likes WHERE user_id = $1 AND target_id = $2 AND type = 'comment'"
    )
    .bind(*user_id)
    .bind(path.into_inner())
    .fetch_one(db.pool())
    .await?;

    if existing > 0 {
        return Err(AppError::Conflict("Already liked".into()));
    }

    sqlx::query("INSERT INTO likes (user_id, target_id, type) VALUES ($1, $2, 'comment')")
        .bind(*user_id)
        .bind(path.into_inner())
        .execute(db.pool())
        .await?;

    sqlx::query("UPDATE comments SET like_count = like_count + 1 WHERE id = $1")
        .bind(path.into_inner())
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().finish())
}

async fn unlike(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let result = sqlx::query(
        "DELETE FROM likes WHERE user_id = $1 AND target_id = $2 AND type = 'comment'"
    )
    .bind(*user_id)
    .bind(path.into_inner())
    .execute(db.pool())
    .await?;

    if result.rows_affected() == 0 {
        return Err(AppError::NotFound("Not liked".into()));
    }

    sqlx::query("UPDATE comments SET like_count = GREATEST(like_count - 1, 0) WHERE id = $1")
        .bind(path.into_inner())
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().finish())
}
