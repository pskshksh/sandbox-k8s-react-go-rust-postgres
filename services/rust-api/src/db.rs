use sqlx::{postgres::PgPoolOptions, PgPool};
use std::{env, time::Duration};

pub async fn connect() -> PgPool {
    let db_url = env::var("DB_URL").expect("DB_URL env variable doesn't exist");

    PgPoolOptions::new()
        .max_connections(25)
        .min_connections(0)
        .acquire_timeout(Duration::from_secs(30))
        .idle_timeout(Some(Duration::from_secs(120)))
        .max_lifetime(Some(Duration::from_secs(300)))
        .connect(&db_url)
        .await
        .expect("could not connect to database")
}

pub async fn migrate(pool: &PgPool) {
    sqlx::query(
        "CREATE TABLE IF NOT EXISTS requests (
            id         SERIAL PRIMARY KEY,
            api_name   TEXT        NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )",
    )
    .execute(pool)
    .await
    .expect("migration failed");

    println!("migration applied successfully");
}
