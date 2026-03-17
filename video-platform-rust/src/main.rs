use actix_web::{middleware, web, App, HttpServer};
use std::net::SocketAddr;
use tracing::info;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod config;
mod db;
mod error;
mod handlers;
mod middleware as app_middleware;
mod models;
mod redis;
mod services;
mod utils;

use config::AppConfig;
use db::Database;
use redis::RedisClient;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenvy::dotenv().ok();

    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new(
            std::env::var("RUST_LOG").unwrap_or_else(|_| "info,video_platform_rust=debug".into()),
        ))
        .with(tracing_subscriber::fmt::layer())
        .init();

    let config = AppConfig::from_env().expect("Failed to load configuration");
    info!("Loaded configuration: {:?}", config.server);

    let db = Database::new(&config.database).await.expect("Failed to connect to database");
    info!("Connected to database");

    let redis_client = RedisClient::new(&config.redis).await.expect("Failed to connect to Redis");
    info!("Connected to Redis");

    let addr: SocketAddr = format!("{}:{}", config.server.host, config.server.port)
        .parse()
        .expect("Invalid address");

    info!("Starting server at {}", addr);

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(db.clone()))
            .app_data(web::Data::new(redis_client.clone()))
            .app_data(web::Data::new(config.clone()))
            .wrap(middleware::Logger::default())
            .wrap(actix_cors::Cors::permissive())
            .wrap(app_middleware::metrics::Metrics)
            .service(
                web::scope("/api/v1")
                    .configure(handlers::auth::configure)
                    .configure(handlers::users::configure)
                    .configure(handlers::videos::configure)
                    .configure(handlers::comments::configure)
                    .configure(handlers::interactions::configure)
                    .configure(handlers::recommend::configure)
                    .configure(handlers::playback::configure)
                    .configure(handlers::danmakus::configure)
                    .configure(handlers::follows::configure)
                    .configure(handlers::search::configure),
            )
            .route("/health", web::get().to(handlers::health))
            .route("/metrics", web::get().to(handlers::metrics))
    })
    .bind(addr)?
    .workers(config.server.workers)
    .run()
    .await
}
