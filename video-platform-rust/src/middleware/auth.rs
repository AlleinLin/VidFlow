use actix_service::{Service, Transform};
use actix_web::{dev::{ServiceRequest, ServiceResponse}, Error, HttpMessage};
use futures::future::LocalBoxFuture;
use std::future::{ready, Ready};

use crate::services::Claims;

pub struct Auth;

impl<S, B> Transform<S, ServiceRequest> for Auth
where
    S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error>,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<B>;
    type Error = Error;
    type InitError = ();
    type Transform = AuthMiddleware<S>;
    type Future = Ready<Result<Self::Transform, Self::InitError>>;

    fn new_transform(&self, service: S) -> Self::Future {
        ready(Ok(AuthMiddleware { service }))
    }
}

pub struct AuthMiddleware<S> {
    service: S,
}

impl<S, B> Service<ServiceRequest> for AuthMiddleware<S>
where
    S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error>,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<B>;
    type Error = Error;
    type Future = LocalBoxFuture<'static, Result<Self::Response, Self::Error>>;

    actix_web::dev::forward_ready!(service);

    fn call(&self, req: ServiceRequest) -> Self::Future {
        let auth_header = req
            .headers()
            .get("Authorization")
            .and_then(|h| h.to_str().ok());

        let fut = self.service.call(req);

        Box::pin(async move {
            match auth_header {
                Some(header) if header.starts_with("Bearer ") => {
                    let token = header[7..].to_string();
                    let config = req.app_data::<actix_web::web::Data<crate::config::AppConfig>>()
                        .map(|c| c.jwt.clone());

                    if let Some(jwt_config) = config {
                        match validate_token(&token, &jwt_config.secret) {
                            Ok(claims) => {
                                req.extensions_mut().insert(claims.user_id);
                            }
                            Err(_) => {
                                return Err(actix_web::error::ErrorUnauthorized("Invalid token"));
                            }
                        }
                    }
                    fut.await
                }
                _ => Err(actix_web::error::ErrorUnauthorized("Missing authorization header")),
            }
        })
    }
}

use jsonwebtoken::{decode, DecodingKey, Validation};

fn validate_token(token: &str, secret: &str) -> Result<Claims, jsonwebtoken::errors::Error> {
    let decoding_key = DecodingKey::from_secret(secret.as_bytes());
    let token_data = decode::<Claims>(token, &decoding_key, &Validation::default())?;

    Ok(token_data.claims)
}
