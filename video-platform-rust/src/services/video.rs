use crate::error::AppError;
use crate::models::*;
use crate::{db::Database, redis::RedisClient};
use std::time::Duration;

pub struct VideoService {
    db: Database,
    redis: RedisClient,
}

impl VideoService {
    pub fn new(db: Database, redis: RedisClient) -> Self {
        Self { db, redis }
    }

    pub async fn create(&self, user_id: i64, input: &CreateVideo) -> Result<Video, AppError> {
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

    pub async fn get_by_id(&self, id: i64) -> Result<Video, AppError> {
        let cache_key = crate::redis::keys::video_meta(id);

        if let Some(cached) = self.redis.get::<Video>(&cache_key).await? {
            return Ok(cached);
        }

        let video = sqlx::query_as::<_, Video>("SELECT * FROM videos WHERE id = $1")
            .bind(id)
            .fetch_optional(self.db.pool())
            .await?
            .ok_or_else(|| AppError::NotFound("Video not found".into()))?;

        let _ = self.redis.set(&cache_key, &video, Duration::from_secs(300)).await;

        Ok(video)
    }

    pub async fn list(
        &self,
        keyword: Option<&str>,
        category_id: Option<i64>,
        page: i32,
        page_size: i32,
    ) -> Result<VideoListResponse, AppError> {
        let offset = (page - 1) * page_size;

        let count_query = if keyword.is_some() {
            "SELECT COUNT(*) FROM videos WHERE status = 'Published' AND (title ILIKE $1 OR description ILIKE $1)"
        } else if category_id.is_some() {
            "SELECT COUNT(*) FROM videos WHERE status = 'Published' AND category_id = $1"
        } else {
            "SELECT COUNT(*) FROM videos WHERE status = 'Published'"
        };

        let total: i64 = match (keyword, category_id) {
            (Some(kw), _) => sqlx::query_scalar(count_query)
                .bind(format!("%{}%", kw))
                .fetch_one(self.db.pool())
                .await?,
            (None, Some(cat)) => sqlx::query_scalar(count_query)
                .bind(cat)
                .fetch_one(self.db.pool())
                .await?,
            (None, None) => sqlx::query_scalar(count_query)
                .fetch_one(self.db.pool())
                .await?,
        };

        let videos: Vec<Video> = match (keyword, category_id) {
            (Some(kw), _) => sqlx::query_as(
                r#"
                SELECT * FROM videos 
                WHERE status = 'Published' AND (title ILIKE $1 OR description ILIKE $1)
                ORDER BY created_at DESC
                LIMIT $2 OFFSET $3
                "#,
            )
            .bind(format!("%{}%", kw))
            .bind(page_size)
            .bind(offset)
            .fetch_all(self.db.pool())
            .await?,
            (None, Some(cat)) => sqlx::query_as(
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
            .await?,
            (None, None) => sqlx::query_as(
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
            .await?,
        };

        Ok(VideoListResponse {
            videos: videos.into_iter().map(VideoResponse::from).collect(),
            total,
            page,
            page_size,
        })
    }

    pub async fn increment_view_count(&self, id: i64) -> Result<(), AppError> {
        sqlx::query("UPDATE videos SET view_count = view_count + 1 WHERE id = $1")
            .bind(id)
            .execute(self.db.pool())
            .await?;

        Ok(())
    }

    pub async fn delete(&self, id: i64, user_id: i64) -> Result<(), AppError> {
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

    pub async fn publish(&self, id: i64, user_id: i64) -> Result<(), AppError> {
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
