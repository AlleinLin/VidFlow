use actix_web::{web, HttpResponse};
use serde::Deserialize;

use crate::db::Database;
use crate::error::AppError;
use crate::models::*;

pub fn configure(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/follows")
            .route("/{user_id}", web::post().to(follow_user))
            .route("/{user_id}", web::delete().to(unfollow_user))
            .route("/followers/{user_id}", web::get().to(get_followers))
            .route("/following/{user_id}", web::get().to(get_following))
    );
}

#[derive(Debug, Deserialize)]
struct PaginationQuery {
    #[serde(default = "default_page")]
    page: i32,
    #[serde(default = "default_page_size")]
    page_size: i32,
}

fn default_page() -> i32 { 1 }
fn default_page_size() -> i32 { 20 }

async fn follow_user(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let following_id = path.into_inner();
    
    if *user_id == following_id {
        return Err(AppError::BadRequest("Cannot follow yourself".into()));
    }

    let existing: i64 = sqlx::query_scalar(
        "SELECT COUNT(*) FROM user_follows WHERE follower_id = $1 AND following_id = $2"
    )
    .bind(*user_id)
    .bind(following_id)
    .fetch_one(db.pool())
    .await?;

    if existing > 0 {
        return Ok(HttpResponse::Ok().json(serde_json::json!({ "message": "Already following" })));
    }

    sqlx::query("INSERT INTO user_follows (follower_id, following_id, created_at) VALUES ($1, $2, NOW())")
        .bind(*user_id)
        .bind(following_id)
        .execute(db.pool())
        .await?;

    sqlx::query("UPDATE users SET following_count = following_count + 1 WHERE id = $1")
        .bind(*user_id)
        .execute(db.pool())
        .await?;

    sqlx::query("UPDATE users SET follower_count = follower_count + 1 WHERE id = $1")
        .bind(following_id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({ "message": "Followed successfully" })))
}

async fn unfollow_user(
    db: web::Data<Database>,
    user_id: web::ReqData<i64>,
    path: web::Path<i64>,
) -> Result<HttpResponse, AppError> {
    let following_id = path.into_inner();

    let result = sqlx::query("DELETE FROM user_follows WHERE follower_id = $1 AND following_id = $2")
        .bind(*user_id)
        .bind(following_id)
        .execute(db.pool())
        .await?;

    if result.rows_affected() == 0 {
        return Ok(HttpResponse::Ok().json(serde_json::json!({ "message": "Not following" })));
    }

    sqlx::query("UPDATE users SET following_count = GREATEST(following_count - 1, 0) WHERE id = $1")
        .bind(*user_id)
        .execute(db.pool())
        .await?;

    sqlx::query("UPDATE users SET follower_count = GREATEST(follower_count - 1, 0) WHERE id = $1")
        .bind(following_id)
        .execute(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({ "message": "Unfollowed successfully" })))
}

#[derive(Debug, Clone, FromRow, Serialize)]
struct PublicUser {
    id: i64,
    username: String,
    display_name: String,
    avatar_url: Option<String>,
    follower_count: i64,
    following_count: i64,
}

async fn get_followers(
    db: web::Data<Database>,
    path: web::Path<i64>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let user_id = path.into_inner();
    let offset = (query.page - 1) * query.page_size;
    
    let users: Vec<PublicUser> = sqlx::query_as(
        r#"
        SELECT u.id, u.username, u.display_name, u.avatar_url, u.follower_count, u.following_count
        FROM users u
        JOIN user_follows uf ON u.id = uf.follower_id
        WHERE uf.following_id = $1
        ORDER BY uf.created_at DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(user_id)
    .bind(query.page_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    let total: i64 = sqlx::query_scalar("SELECT COUNT(*) FROM user_follows WHERE following_id = $1")
        .bind(user_id)
        .fetch_one(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "users": users,
        "total": total,
        "page": query.page,
        "page_size": query.page_size
    })))
}

async fn get_following(
    db: web::Data<Database>,
    path: web::Path<i64>,
    query: web::Query<PaginationQuery>,
) -> Result<HttpResponse, AppError> {
    let user_id = path.into_inner();
    let offset = (query.page - 1) * query.page_size;
    
    let users: Vec<PublicUser> = sqlx::query_as(
        r#"
        SELECT u.id, u.username, u.display_name, u.avatar_url, u.follower_count, u.following_count
        FROM users u
        JOIN user_follows uf ON u.id = uf.following_id
        WHERE uf.follower_id = $1
        ORDER BY uf.created_at DESC
        LIMIT $2 OFFSET $3
        "#,
    )
    .bind(user_id)
    .bind(query.page_size)
    .bind(offset)
    .fetch_all(db.pool())
    .await?;

    let total: i64 = sqlx::query_scalar("SELECT COUNT(*) FROM user_follows WHERE follower_id = $1")
        .bind(user_id)
        .fetch_one(db.pool())
        .await?;

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "users": users,
        "total": total,
        "page": query.page,
        "page_size": query.page_size
    })))
}
