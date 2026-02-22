package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

type RequestsResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"count"`
}

type InsertRequestBody struct {
	Name string `json:"name"`
}

func (h *Handler) insertRequest(w http.ResponseWriter, r *http.Request) {
	var body InsertRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Name != "go" && body.Name != "rust" {
		respondError(w, http.StatusBadRequest, "name must be 'go' or 'rust'")
		return
	}

	_, err := h.db.Exec(`INSERT INTO requests (api_name) VALUES ($1)`, body.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) liveness(w http.ResponseWriter, r *http.Request) {
	err := h.db.Ping()
	if err != nil {
		respondError(w, http.StatusServiceUnavailable, "db ping failed")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) readiness(w http.ResponseWriter, r *http.Request) {
	_, err := h.db.Exec("SELECT 1")
	if err != nil {
		respondError(w, http.StatusServiceUnavailable, "db not ready")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (h *Handler) getRequests(w http.ResponseWriter, r *http.Request) {
	var resp RequestsResponse
	err := h.db.QueryRow(`SELECT NOW(), COUNT(*) FROM requests WHERE api_name = 'go'`).Scan(&resp.Timestamp, &resp.Count)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, resp)
}
