use axum::{
    extract::State,
    http::StatusCode,
    response::Json,
};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use serde_json::{json, Value};
use sqlx::{FromRow, PgPool};

type JsonError = (StatusCode, Json<Value>);

fn internal_error(e: impl ToString) -> JsonError {
    (
        StatusCode::INTERNAL_SERVER_ERROR,
        Json(json!({"error": e.to_string()})),
    )
}

#[derive(Serialize, FromRow)]
pub struct RequestsResponse {
    pub timestamp: DateTime<Utc>,
    pub count: i64,
}

#[derive(Deserialize)]
pub struct InsertRequestBody {
    pub name: String,
}

pub async fn liveness(State(pool): State<PgPool>) -> Result<Json<Value>, JsonError> {
    sqlx::query("SELECT 1")
        .execute(&pool)
        .await
        .map_err(|_| {
            (
                StatusCode::SERVICE_UNAVAILABLE,
                Json(json!({"error": "db ping failed"})),
            )
        })?;
    Ok(Json(json!({"status": "ok"})))
}

pub async fn readiness(State(pool): State<PgPool>) -> Result<Json<Value>, JsonError> {
    sqlx::query("SELECT 1")
        .execute(&pool)
        .await
        .map_err(|_| {
            (
                StatusCode::SERVICE_UNAVAILABLE,
                Json(json!({"error": "db not ready"})),
            )
        })?;
    Ok(Json(json!({"status": "ready"})))
}

pub async fn get_requests(
    State(pool): State<PgPool>,
) -> Result<Json<RequestsResponse>, JsonError> {
    let row = sqlx::query_as::<_, RequestsResponse>(
        "SELECT NOW() as timestamp, COUNT(*) as count FROM requests WHERE api_name = 'rust'",
    )
    .fetch_one(&pool)
    .await
    .map_err(internal_error)?;

    Ok(Json(row))
}

pub async fn insert_request(
    State(pool): State<PgPool>,
    Json(body): Json<InsertRequestBody>,
) -> Result<StatusCode, JsonError> {
    if body.name != "rust" {
        return Err((
            StatusCode::BAD_REQUEST,
            Json(json!({"error": "name must be 'go' or 'rust'"})),
        ));
    }

    sqlx::query("INSERT INTO requests (api_name) VALUES ($1)")
        .bind(&body.name)
        .execute(&pool)
        .await
        .map_err(internal_error)?;

    Ok(StatusCode::CREATED)
}
