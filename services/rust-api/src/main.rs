mod db;
mod handler;

use axum::{
    routing::{get, post},
    Router,
};

#[tokio::main]
async fn main() {
    let pool = db::connect().await;
    db::migrate(&pool).await;

    let app = Router::new()
        .route("/healthz", get(handler::liveness))
        .route("/readyz", get(handler::readiness))
        .route("/requests", get(handler::get_requests))
        .route("/requests", post(handler::insert_request))
        .with_state(pool);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:8081").await.unwrap();
    println!("server started on :8081");
    axum::serve(listener, app).await.unwrap();
}
