use actix_web::{http::StatusCode, HttpResponse, ResponseError};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("Not found: {0}")]
    NotFound(String),

    #[error("Unauthorized: {0}")]
    Unauthorized(String),

    #[error("BadRequest: {0}")]
    BadRequest(String),

    #[error("Conflict: {0}")]
    Conflict(String),

    #[error("Internal error: {0}")]
    Internal(String),

    #[error("Database error: {0}")]
    DatabaseError(#[from] sqlx::Error),

    #[error("Redis error: {0}")]
    RedisError(#[from] redis::RedisError),

    #[error("Validation error: {0}")]
    ValidationError(String),
}

impl ResponseError for AppError {
    fn error_response(&self) -> HttpResponse {
        let status_code = self.status_code();
        let error_message = self.to_string();

        HttpResponse::build(status_code).json(serde_json::json!({
            "error": error_message,
            "code": status_code.as_u16(),
        }))
    }
}

impl AppError {
    fn status_code(&self) -> StatusCode {
        match self {
            AppError::NotFound(_) => StatusCode::NOT_FOUND,
            AppError::Unauthorized(_) => StatusCode::UNAUTHORIZED,
            AppError::BadRequest(_) => StatusCode::BAD_REQUEST,
            AppError::Conflict(_) => StatusCode::CONFLICT,
            AppError::ValidationError(_) => StatusCode::BAD_REQUEST,
            AppError::Internal(_) | AppError::DatabaseError(_) | AppError::RedisError(_) => {
                StatusCode::INTERNAL_SERVER_ERROR
            }
        }
    }
}

pub type AppResult<T> = Result<T, AppError>;
