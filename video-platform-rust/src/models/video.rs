use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "VARCHAR")]
pub enum VideoStatus {
    Uploading,
    Transcoding,
    Auditing,
    Published,
    Rejected,
    Deleted,
    Hidden,
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "VARCHAR")]
pub enum VideoVisibility {
    Public,
    FollowersOnly,
    Private,
}

impl Default for VideoVisibility {
    fn default() -> Self {
        Self::Public
    }
}

#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct Video {
    pub id: i64,
    pub user_id: i64,
    pub title: String,
    pub description: Option<String>,
    pub status: VideoStatus,
    pub visibility: VideoVisibility,
    pub duration_seconds: Option<i32>,
    pub original_filename: Option<String>,
    pub storage_key: Option<String>,
    pub thumbnail_url: Option<String>,
    pub category_id: Option<i64>,
    pub view_count: i64,
    pub like_count: i64,
    pub comment_count: i32,
    pub published_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CreateVideo {
    pub title: String,
    pub description: Option<String>,
    pub category_id: Option<i64>,
    pub visibility: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VideoResponse {
    pub id: i64,
    pub user_id: i64,
    pub title: String,
    pub description: Option<String>,
    pub status: String,
    pub visibility: String,
    pub duration_seconds: Option<i32>,
    pub thumbnail_url: Option<String>,
    pub category_id: Option<i64>,
    pub view_count: i64,
    pub like_count: i64,
    pub comment_count: i32,
    pub published_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
}

impl From<Video> for VideoResponse {
    fn from(video: Video) -> Self {
        Self {
            id: video.id,
            user_id: video.user_id,
            title: video.title,
            description: video.description,
            status: match video.status {
                VideoStatus::Uploading => "UPLOADING".to_string(),
                VideoStatus::Transcoding => "TRANSCODING".to_string(),
                VideoStatus::Auditing => "AUDITING".to_string(),
                VideoStatus::Published => "PUBLISHED".to_string(),
                VideoStatus::Rejected => "REJECTED".to_string(),
                VideoStatus::Deleted => "DELETED".to_string(),
                VideoStatus::Hidden => "HIDDEN".to_string(),
            },
            visibility: match video.visibility {
                VideoVisibility::Public => "PUBLIC".to_string(),
                VideoVisibility::FollowersOnly => "FOLLOWERS_ONLY".to_string(),
                VideoVisibility::Private => "PRIVATE".to_string(),
            },
            duration_seconds: video.duration_seconds,
            thumbnail_url: video.thumbnail_url,
            category_id: video.category_id,
            view_count: video.view_count,
            like_count: video.like_count,
            comment_count: video.comment_count,
            published_at: video.published_at,
            created_at: video.created_at,
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VideoListResponse {
    pub videos: Vec<VideoResponse>,
    pub total: i64,
    pub page: i32,
    pub page_size: i32,
}
