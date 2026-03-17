use serde::Deserialize;
use std::env;

#[derive(Debug, Clone, Deserialize)]
pub struct AppConfig {
    pub server: ServerConfig,
    pub database: DatabaseConfig,
    pub redis: RedisConfig,
    pub jwt: JwtConfig,
    pub kafka: KafkaConfig,
    pub storage: StorageConfig,
}

#[derive(Debug, Clone, Deserialize)]
pub struct ServerConfig {
    pub host: String,
    pub port: u16,
    pub workers: usize,
}

#[derive(Debug, Clone, Deserialize)]
pub struct DatabaseConfig {
    pub url: String,
    pub max_connections: u32,
    pub min_connections: u32,
}

#[derive(Debug, Clone, Deserialize)]
pub struct RedisConfig {
    pub url: String,
    pub pool_size: u32,
}

#[derive(Debug, Clone, Deserialize)]
pub struct JwtConfig {
    pub secret: String,
    pub access_token_ttl: i64,
    pub refresh_token_ttl: i64,
    pub issuer: String,
}

#[derive(Debug, Clone, Deserialize)]
pub struct KafkaConfig {
    pub brokers: Vec<String>,
    pub topic: String,
    pub consumer_group: String,
}

#[derive(Debug, Clone, Deserialize)]
pub struct StorageConfig {
    pub endpoint: String,
    pub access_key: String,
    pub secret_key: String,
    pub bucket: String,
    pub region: String,
}

impl AppConfig {
    pub fn from_env() -> Result<Self, config::ConfigError> {
        config::Config::builder()
            .set_default("server.host", "0.0.0.0")?
            .set_default("server.port", 8082)?
            .set_default("server.workers", 4)?
            .set_default("database.max_connections", 100)?
            .set_default("database.min_connections", 10)?
            .set_default("redis.pool_size", 100)?
            .set_default("jwt.access_token_ttl", 900)?
            .set_default("jwt.refresh_token_ttl", 604800)?
            .set_default("jwt.issuer", "video-platform")?
            .set_default("kafka.topic", "video.events")?
            .set_default("storage.region", "us-east-1")?
            .add_source(config::Environment::default().separator("__"))
            .build()?
            .try_deserialize()
    }
}
