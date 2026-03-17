use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "VARCHAR")]
pub enum UserRole {
    User,
    Creator,
    Moderator,
    Admin,
}

impl Default for UserRole {
    fn default() -> Self {
        Self::User
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "VARCHAR")]
pub enum UserStatus {
    Active,
    Suspended,
    Banned,
}

impl Default for UserStatus {
    fn default() -> Self {
        Self::Active
    }
}

#[derive(Debug, Clone, FromRow, Serialize, Deserialize)]
pub struct User {
    pub id: i64,
    pub username: String,
    pub email: String,
    #[serde(skip_serializing)]
    pub password_hash: String,
    pub display_name: String,
    pub avatar_url: Option<String>,
    pub bio: Option<String>,
    pub role: UserRole,
    pub status: UserStatus,
    pub follower_count: i64,
    pub following_count: i64,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub last_login_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CreateUser {
    pub username: String,
    pub email: String,
    pub password: String,
    pub display_name: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UserProfile {
    pub id: i64,
    pub username: String,
    pub email: String,
    pub display_name: String,
    pub avatar_url: Option<String>,
    pub bio: Option<String>,
    pub role: String,
    pub status: String,
    pub follower_count: i64,
    pub following_count: i64,
    pub created_at: DateTime<Utc>,
}

impl From<User> for UserProfile {
    fn from(user: User) -> Self {
        Self {
            id: user.id,
            username: user.username,
            email: user.email,
            display_name: user.display_name,
            avatar_url: user.avatar_url,
            bio: user.bio,
            role: match user.role {
                UserRole::User => "USER".to_string(),
                UserRole::Creator => "CREATOR".to_string(),
                UserRole::Moderator => "MODERATOR".to_string(),
                UserRole::Admin => "ADMIN".to_string(),
            },
            status: match user.status {
                UserStatus::Active => "ACTIVE".to_string(),
                UserStatus::Suspended => "SUSPENDED".to_string(),
                UserStatus::Banned => "BANNED".to_string(),
            },
            follower_count: user.follower_count,
            following_count: user.following_count,
            created_at: user.created_at,
        }
    }
}
