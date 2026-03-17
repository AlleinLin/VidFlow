use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;

#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct WatchHistory {
    pub id: i64,
    pub user_id: i64,
    pub video_id: i64,
    pub watch_duration: i64,
    pub watch_progress: f64,
    pub last_position: f64,
    pub completed: bool,
    pub watched_at: DateTime<Utc>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UpdateProgressRequest {
    pub video_id: i64,
    pub position: f64,
    pub duration: f64,
    pub watch_duration: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WatchHistoryResponse {
    pub id: i64,
    pub video_id: i64,
    pub video_title: String,
    pub video_thumbnail: Option<String>,
    pub watch_progress: f64,
    pub last_position: f64,
    pub completed: bool,
    pub watched_at: DateTime<Utc>,
}

#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct Danmaku {
    pub id: i64,
    pub video_id: i64,
    pub user_id: i64,
    pub content: String,
    pub position_seconds: f64,
    pub style: String,
    pub color: String,
    pub font_size: i32,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CreateDanmakuRequest {
    pub video_id: i64,
    pub content: String,
    pub position_seconds: f64,
    #[serde(default)]
    pub style: Option<String>,
    #[serde(default = "default_color")]
    pub color: String,
    #[serde(default = "default_font_size")]
    pub font_size: i32,
}

fn default_color() -> String { "#FFFFFF".to_string() }
fn default_font_size() -> i32 { 24 }

#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct UserFollow {
    pub id: i64,
    pub follower_id: i64,
    pub following_id: i64,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct FollowUser {
    pub following_id: i64,
}
