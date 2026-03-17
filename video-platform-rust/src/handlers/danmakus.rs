use actix_web::{web, HttpResponse};
use serde::Deserialize;

use crate::db::Database;
use crate::error::AppError;
use crate::models::*;

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/danmakus")
            .route("/video/{video_id}", web::get().to(get_danmakus))
            .route("/", web::post().to(create_danmaku))
            .route("/{id}", web::delete().to(delete_danmaku))
    );
}

#[derive(Debug, Deserialize)]
struct PositionQuery {
    #[serde(default)]
    start: Option<f64>,
    #[serde(default)]
    end: Option<f64>,
}

async fn get_danmakus(
    db: web::Data<Database>,
    path: web::Path<i64>,
    query: web::Query<PositionQuery>,
) -> Result<HttpResponse, AppError> {
    let video_id = path.into_inner();
    
    let danmakus: Vec<Danmaku> = match (query.start, query.end) {
        (Some(start), Some(end)) => {
            sqlx::query_as(
                "SELECT id, video_id, user_id, content, position_seconds, style, color, font_size, created_at FROM danmakus WHERE video_id = $1 AND position_seconds >= $2 AND position_seconds <= $3 ORDER BY position_seconds"
            )
            .bind(video_id)
            .bind(start)
            .bind(end)
            .fetch_all(db.pool())
            .await?
        },
        _ => {
            sqlx::query_as(
                "SELECT id, video_id, user_id, content, position_seconds, style, color, font_size, created_at FROM danmakus WHERE video_id = $1 ORDER BY position_seconds"
            )
            .bind(video_id)
            .fetch_all(db.pool())
            .await?
        }
    };

    Ok(HttpResponse::Ok().json(serde_json::json!({ "danmakus": danmakus })))
}

async fn create_danmaku(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    req: web::Json<CreateDanmakuRequest>,
) -> Result<HttpResponse, AppError> {
    let style = req.style.clone().unwrap_or_else(|| "SCROLL".to_string());
    
    let danmaku: Danmaku = sqlx::query_as(
        r#"
        INSERT INTO danmakus (video_id, user_id, content, position_seconds, style, color, font_size, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
        RETURNING id, video_id, user_id, content, position_seconds, style, color, font_size, created_at
        "#,
    )
    .bind(req.video_id)
    .bind(*user_id)
    .bind(&req.content)
    .bind(req.position_seconds)
    .bind(&style)
    .bind(&req.color)
    .bind(req.font_size)
    .fetch_one(db.pool())
    .await?;

    Ok(HttpResponse::Created().json(danmaku))
}

async fn delete_danmaku(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let danmaku_id = path.into_inner();
    
    let result = sqlx::query("DELETE FROM danmakus WHERE id = $1 AND user_id = $2")
        .bind(danmaku_id)
        .bind(*user_id)
        .execute(db.pool())
        .await?;

    if result.rows_affected() == 0 {
        return Err(AppError::NotFound("Danmaku not found".into()));
    }

    Ok(HttpResponse::NoContent().finish())
}
