use actix_web::{web, HttpResponse};
use serde::Deserialize;

use crate::db::Database;
use crate::error::AppError;
use crate::models::*;
use crate::redis::RedisClient;

#[derive(Debug, serde::Deserialize)]
pub struct CreateVideoRequest {
    pub title: String,
    pub description: Option<String>,
    pub category_id: Option<i64>,
    pub visibility: Option<String>,
}

#[derive(Debug, serde::Deserialize)]
pub struct VideoQuery {
    pub keyword: Option<String>,
    pub category_id: Option<i64>,
    #[serde(default = "default_page")]
    pub page: i32,
    #[serde(default = "default_page_size")]
    pub page_size: i32,
}

fn default_page() -> i32 { 1 }
fn default_page_size() -> i32 { 20 }

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/videos")
            .route("", web::get().to(list))
            .route("", web::post().to(create))
            .route("/{id}", web::get().to(get))
            .route("/{id}", web::put().to(update))
            .route("/{id}", web::delete().to(delete))
            .route("/{id}/publish", web::post().to(publish)),
    );
}

async fn list(
    db: web::Data<Database>,
    redis: web::Data<RedisClient>,
    query: web::Query<VideoQuery>,
) -> Result<HttpResponse, AppError> {
    let service = VideoService::new(db.into_inner(), redis.into_inner());
    let response = service
        .list(
            query.keyword.as_deref(),
            query.category_id,
            query.page,
            query.page_size,
        )
        .await?;

    Ok(HttpResponse::Ok().json(response))
}

async fn get(
    db: web::Data<Database>,
    redis: web::Data<RedisClient>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let service = VideoService::new(db.into_inner(), redis.into_inner());
    let video = service.get_by_id(path.into_inner()).await?;

    let _ = service.increment_view_count(video.id).await;

    Ok(HttpResponse::Ok().json(VideoResponse::from(video)))
}

async fn create(
    db: web::Data<Database>,
    redis: web::Data<RedisClient>,
    user_id: web::ReqData<i64>,
    body: web::Json<CreateVideoRequest>,
) -> Result<HttpResponse, AppError> {
    let service = VideoService::new(db.into_inner(), redis.into_inner());
    let input = CreateVideo {
        title: body.title.clone(),
        description: body.description.clone(),
        category_id: body.category_id,
        visibility: body.visibility.clone(),
    };

    let video = service.create(*user_id, &input).await?;

    Ok(HttpResponse::Ok().json(VideoResponse::from(video)))
}

async fn update(
    db: web::Data<Database>,
    redis: web::Data<RedisClient>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
    body: web::Json<CreateVideoRequest>,
) -> Result<HttpResponse, AppError> {
    let service = VideoService::new(db.into_inner(), redis.into_inner());
    let video = service.update(
        path.into_inner(),
        *user_id,
        &body.title,
        body.description.as_deref(),
        body.visibility.as_deref(),
    ).await?;

    Ok(HttpResponse::Ok().json(VideoResponse::from(video)))
}

async fn delete(
    db: web::Data<Database>,
    redis: web::Data<RedisClient>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let service = VideoService::new(db.into_inner(), redis.into_inner());
    service.delete(path.into_inner(), *user_id).await?;

    Ok(HttpResponse::Ok().finish())
}

async fn publish(
    db: web::Data<Database>,
    redis: web::Data<RedisClient>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let service = VideoService::new(db.into_inner(), redis.into_inner());
    service.publish(path.into_inner(), *user_id).await?;

    Ok(HttpResponse::Ok().finish())
}

struct VideoService {
    db: Database,
    redis: RedisClient,
}

impl VideoService {
    fn new(db: Database, redis: RedisClient) -> Self {
        Self { db, redis }
    }

    async fn create(&self, user_id: i64, input: &CreateVideo) -> Result<Video, AppError> {
        let visibility = input
            .visibility
            .as_ref()
            .map(|v| match v.as_str() {
                "PUBLIC" => VideoVisibility::Public,
                "FOLLOWERS_ONLY" => VideoVisibility::FollowersOnly,
                "PRIVATE" => VideoVisibility::Private,
                _ => VideoVisibility::Public,
            })
            .unwrap_or_default();

        let video = sqlx::query_as::<_, Video>(
            r#"
            INSERT INTO videos (user_id, title, description, status, visibility, category_id)
            VALUES ($1, $2, $3, 'Uploading', $4, $5)
            RETURNING *
            "#,
        )
        .bind(user_id)
        .bind(&input.title)
        .bind(&input.description)
        .bind(&visibility)
        .bind(input.category_id)
        .fetch_one(self.db.pool())
        .await?;

        Ok(video)
    }

    async fn get_by_id(&self, id: i64) -> Result<Video, AppError> {
        let cache_key = crate::redis::keys::video_meta(id);

        if let Some(cached) = self.redis.get::<Video>(&cache_key).await? {
            return Ok(cached);
        }

        let video = sqlx::query_as::<_, Video>("SELECT * FROM videos WHERE id = $1")
            .bind(id)
            .fetch_optional(self.db.pool())
            .await?
            .ok_or_else(|| AppError::NotFound("Video not found".into()))?;

        let _ = self.redis.set(&cache_key, &video, std::time::Duration::from_secs(300)).await;

        Ok(video)
    }

    async fn update(
        &self,
        id: i64,
        user_id: i64,
        title: &str,
        description: Option<&str>,
        visibility: Option<&str>,
    ) -> Result<Video, AppError> {
        let video = self.get_by_id(id).await?;
        
        if video.user_id != user_id {
            return Err(AppError::Unauthorized("Unauthorized".into()));
        }

        let visibility_enum = visibility.map(|v| match v {
            "PUBLIC" => VideoVisibility::Public,
            "FOLLOWERS_ONLY" => VideoVisibility::FollowersOnly,
            "PRIVATE" => VideoVisibility::Private,
            _ => video.visibility.clone(),
        });

        let updated = sqlx::query_as::<_, Video>(
            r#"
            UPDATE videos 
            SET title = $1, description = $2, visibility = COALESCE($3, visibility)
            WHERE id = $4
            RETURNING *
            "#,
        )
        .bind(title)
        .bind(description)
        .bind(visibility_enum)
        .bind(id)
        .fetch_one(self.db.pool())
        .await?;

        let _ = self.redis.del(&crate::redis::keys::video_meta(id)).await;

        Ok(updated)
    }

    async fn list(
        &self,
        keyword: Option<&str>,
        category_id: Option<i64>,
        page: i32,
        page_size: i32,
    ) -> Result<VideoListResponse, AppError> {
        let offset = (page - 1) * page_size;

        let (total, videos): (i64, Vec<Video>) = match (keyword, category_id) {
            (Some(kw), _) => {
                let pattern = format!("%{}%", kw);
                let total: i64 = sqlx::query_scalar(
                    "SELECT COUNT(*) FROM videos WHERE status = 'Published' AND (title ILIKE $1 OR description ILIKE $1)"
                )
                .bind(&pattern)
                .fetch_one(self.db.pool())
                .await?;

                let videos: Vec<Video> = sqlx::query_as(
                    r#"
                    SELECT * FROM videos 
                    WHERE status = 'Published' AND (title ILIKE $1 OR description ILIKE $1)
                    ORDER BY created_at DESC
                    LIMIT $2 OFFSET $3
                    "#,
                )
                .bind(&pattern)
                .bind(page_size)
                .bind(offset)
                .fetch_all(self.db.pool())
                .await?;

                (total, videos)
            }
            (None, Some(cat)) => {
                let total: i64 = sqlx::query_scalar(
                    "SELECT COUNT(*) FROM videos WHERE status = 'Published' AND category_id = $1"
                )
                .bind(cat)
                .fetch_one(self.db.pool())
                .await?;

                let videos: Vec<Video> = sqlx::query_as(
                    r#"
                    SELECT * FROM videos 
                    WHERE status = 'Published' AND category_id = $1
                    ORDER BY created_at DESC
                    LIMIT $2 OFFSET $3
                    "#,
                )
                .bind(cat)
                .bind(page_size)
                .bind(offset)
                .fetch_all(self.db.pool())
                .await?;

                (total, videos)
            }
            (None, None) => {
                let total: i64 = sqlx::query_scalar(
                    "SELECT COUNT(*) FROM videos WHERE status = 'Published'"
                )
                .fetch_one(self.db.pool())
                .await?;

                let videos: Vec<Video> = sqlx::query_as(
                    r#"
                    SELECT * FROM videos 
                    WHERE status = 'Published'
                    ORDER BY created_at DESC
                    LIMIT $1 OFFSET $2
                    "#,
                )
                .bind(page_size)
                .bind(offset)
                .fetch_all(self.db.pool())
                .await?;

                (total, videos)
            }
        };

        Ok(VideoListResponse {
            videos: videos.into_iter().map(VideoResponse::from).collect(),
            total,
            page,
            page_size,
        })
    }

    async fn increment_view_count(&self, id: i64) -> Result<(), AppError> {
        sqlx::query("UPDATE videos SET view_count = view_count + 1 WHERE id = $1")
            .bind(id)
            .execute(self.db.pool())
            .await?;

        Ok(())
    }

    async fn delete(&self, id: i64, user_id: i64) -> Result<(), AppError> {
        let result = sqlx::query(
            "UPDATE videos SET status = 'Deleted' WHERE id = $1 AND user_id = $2",
        )
        .bind(id)
        .bind(user_id)
        .execute(self.db.pool())
        .await?;

        if result.rows_affected() == 0 {
            return Err(AppError::NotFound("Video not found or unauthorized".into()));
        }

        let _ = self.redis.del(&crate::redis::keys::video_meta(id)).await;

        Ok(())
    }

    async fn publish(&self, id: i64, user_id: i64) -> Result<(), AppError> {
        let result = sqlx::query(
            r#"
            UPDATE videos 
            SET status = 'Published', published_at = NOW() 
            WHERE id = $1 AND user_id = $2 AND status IN ('Auditing', 'Uploading')
            "#,
        )
        .bind(id)
        .bind(user_id)
        .execute(self.db.pool())
        .await?;

        if result.rows_affected() == 0 {
            return Err(AppError::BadRequest(
                "Video cannot be published in current status".into(),
            ));
        }

        let _ = self.redis.del(&crate::redis::keys::video_meta(id)).await;

        Ok(())
    }
}
