use actix_service::{Service, Transform};
use actix_web::{dev::{ServiceRequest, ServiceResponse}, Error};
use futures::future::LocalBoxFuture;
use std::future::Ready;
use std::time::Instant;

use super::HTTP_REQUESTS_TOTAL;
use super::HTTP_REQUEST_DURATION;

pub struct Metrics;

impl<S, B> Transform<S, ServiceRequest> for Metrics
where
    S: actix_service::Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error>,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<B>;
    type Error = Error;
    type InitError = ();
    type Transform = MetricsMiddleware<S>;
    type Future = Ready<Result<Self::Transform, Self::InitError>>;

    fn new_transform(&self, service: S) -> Self::Future {
        std::future::ready(Ok(MetricsMiddleware { service }))
    }
}

pub struct MetricsMiddleware<S> {
    service: S,
}

impl<S, B> actix_service::Service<ServiceRequest> for MetricsMiddleware<S>
where
    S: actix_service::Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error>,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<B>;
    type Error = Error;
    type Future = LocalBoxFuture<'static, Result<Self::Response, Self::Error>>;

    actix_web::dev::forward_ready!(service);

    fn call(&self, req: ServiceRequest) -> Self::Future {
        let start = Instant::now();
        let method = req.method().to_string();
        let path = req.path().to_string();

        let fut = self.service.call(req);

        Box::pin(async move {
            let res = fut.await;
            let duration = start.elapsed().as_secs_f64();

            let status = match &res {
                Ok(response) => response.status().as_u16().to_string(),
                Err(_) => "500".to_string(),
            };

            HTTP_REQUESTS_TOTAL
                .with_label_values(&[&method, &path, &status])
                .inc();

            HTTP_REQUEST_DURATION
                .with_label_values(&[&method, &path])
                .observe(duration);

            res
        })
    }
}
