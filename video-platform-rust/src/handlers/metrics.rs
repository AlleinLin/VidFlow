use actix_web::HttpResponse;
use prometheus::{Encoder, TextEncoder};

lazy_static::lazy_static! {
    pub static ref HTTP_REQUESTS_TOTAL: prometheus::CounterVec = prometheus::register_counter_vec!(
        "http_requests_total",
        "Total number of HTTP requests",
        &["method", "path", "status"]
    ).unwrap();

    pub static ref HTTP_REQUEST_DURATION: prometheus::HistogramVec = prometheus::register_histogram_vec!(
        "http_request_duration_seconds",
        "Duration of HTTP requests in seconds",
        &["method", "path"]
    ).unwrap();
}

pub async fn metrics() -> HttpResponse {
    let encoder = TextEncoder::new();
    let metric_families = prometheus::gather();
    let mut buffer = Vec::new();

    if let Err(e) = encoder.encode(&metric_families, &mut buffer) {
        return HttpResponse::InternalServerError().body(format!("Failed to encode metrics: {}", e));
    }

    HttpResponse::Ok()
        .content_type("text/plain; version=0.0.4")
        .body(buffer)
}
