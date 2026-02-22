package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const (
	DB_URL = "DB_URL"
)

var DB *sql.DB

func getEnv(envKey, fallback string) string {
	envVal := os.Getenv(envKey)
	if envVal == "" {
		return fallback
	}
	return envVal
}

func migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS requests (
		id         SERIAL PRIMARY KEY,
		api_name       TEXT        NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	return err
}

func init() {
	dbUrl := getEnv(DB_URL, "")
	if dbUrl == "" {
		log.Fatal("env variable db url doesn't exist")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("could not open the db err: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatalf("db ping failed err: %v", err)
	}

	err = migrate(db)
	if err != nil {
		log.Fatalf("migration failed err: %v", err)
	}

	DB = db

	log.Println("migration applied successfully")

}
