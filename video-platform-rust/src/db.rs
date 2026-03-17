use sqlx::postgres::{PgPool, PgPoolOptions};
use std::time::Duration;

#[derive(Clone)]
pub struct Database {
    pool: PgPool,
}

impl Database {
    pub async fn new(config: &super::config::DatabaseConfig) -> Result<Self, sqlx::Error> {
        let pool = PgPoolOptions::new()
            .max_connections(config.max_connections)
            .min_connections(config.min_connections)
            .acquire_timeout(Duration::from_secs(5))
            .connect(&config.url)
            .await?;

        sqlx::migrate!("./migrations")
            .run(&pool)
            .await
            .ok();

        Ok(Self { pool })
    }

    pub fn pool(&self) -> &PgPool {
        &self.pool
    }
}

impl std::ops::Deref for Database {
    type Target = PgPool;

    fn deref(&self) -> &Self::Target {
        &self.pool
    }
}
