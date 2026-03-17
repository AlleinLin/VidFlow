use redis::{Client, ConnectionManager};
use serde::{de::DeserializeOwned, Serialize};
use std::time::Duration;

#[derive(Clone)]
pub struct RedisClient {
    manager: ConnectionManager,
}

impl RedisClient {
    pub async fn new(config: &super::config::RedisConfig) -> Result<Self, redis::RedisError> {
        let client = Client::open(config.url.as_str())?;
        let manager = ConnectionManager::new(client).await?;
        Ok(Self { manager })
    }

    pub async fn get<T: DeserializeOwned>(&self, key: &str) -> Result<Option<T>, redis::RedisError> {
        let mut conn = self.manager.clone();
        let result: Option<String> = redis::cmd("GET")
            .arg(key)
            .query_async(&mut conn)
            .await?;

        match result {
            Some(json) => {
                let value: T = serde_json::from_str(&json).map_err(|_| {
                    redis::RedisError::from((
                        redis::ErrorKind::TypeError,
                        "Failed to deserialize",
                    ))
                })?;
                Ok(Some(value))
            }
            None => Ok(None),
        }
    }

    pub async fn set<T: Serialize>(&self, key: &str, value: &T, ttl: Duration) -> Result<(), redis::RedisError> {
        let json = serde_json::to_string(value).map_err(|_| {
            redis::RedisError::from((
                redis::ErrorKind::TypeError,
                "Failed to serialize",
            ))
        })?;

        let mut conn = self.manager.clone();
        redis::cmd("SET")
            .arg(key)
            .arg(&json)
            .arg("EX")
            .arg(ttl.as_secs())
            .query_async(&mut conn)
            .await
    }

    pub async fn del(&self, key: &str) -> Result<(), redis::RedisError> {
        let mut conn = self.manager.clone();
        redis::cmd("DEL")
            .arg(key)
            .query_async(&mut conn)
            .await
    }

    pub async fn incr(&self, key: &str) -> Result<i64, redis::RedisError> {
        let mut conn = self.manager.clone();
        redis::cmd("INCR")
            .arg(key)
            .query_async(&mut conn)
            .await
    }

    pub async fn sadd(&self, key: &str, member: &str) -> Result<(), redis::RedisError> {
        let mut conn = self.manager.clone();
        redis::cmd("SADD")
            .arg(key)
            .arg(member)
            .query_async(&mut conn)
            .await
    }

    pub async fn srem(&self, key: &str, member: &str) -> Result<(), redis::RedisError> {
        let mut conn = self.manager.clone();
        redis::cmd("SREM")
            .arg(key)
            .arg(member)
            .query_async(&mut conn)
            .await
    }

    pub async fn sismember(&self, key: &str, member: &str) -> Result<bool, redis::RedisError> {
        let mut conn = self.manager.clone();
        redis::cmd("SISMEMBER")
            .arg(key)
            .arg(member)
            .query_async(&mut conn)
            .await
    }
}

pub mod keys {
    pub fn video_meta(video_id: i64) -> String {
        format!("video:meta:{}", video_id)
    }

    pub fn user_profile(user_id: i64) -> String {
        format!("user:profile:{}", user_id)
    }

    pub fn view_count(video_id: i64) -> String {
        format!("stats:view:{}", video_id)
    }

    pub fn like_count(video_id: i64) -> String {
        format!("stats:like:{}", video_id)
    }

    pub fn user_like_set(user_id: i64) -> String {
        format!("user:likes:{}", user_id)
    }

    pub fn recommend_home(user_id: i64) -> String {
        format!("recommend:home:{}", user_id)
    }
}
