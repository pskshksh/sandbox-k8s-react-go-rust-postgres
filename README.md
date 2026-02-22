# sandbox-k8s-react-go-rust-postgres

A sandbox project running two REST APIs written in — **Go** and **Rust** — both backed by the same **PostgreSQL** database. A **React** frontend displays live call counts for each API. The entire stack is containerised and deployed on **Kubernetes**.

---

## What it does

Each API exposes the same contract: you hit a button in the UI, a POST request is recorded in Postgres, and the frontend immediately reflects the updated count. Both APIs share the same `requests` table, so you can compare how many calls each language has handled side by side.

```
  Browser
     │
     ▼
 [React Frontend]   ← shows call counts for Go and Rust, buttons to trigger calls
     │         │
 [Go API]   [Rust API]   ← each records calls to Postgres and serves its own count
     │         │
    [PostgreSQL]          ← single shared database, one table: requests
```

---

## Stack

| Layer      | Technology                              |
|------------|-----------------------------------------|
| Frontend   | React + Vite, served by nginx           |
| Go API     | Go, gorilla/mux, lib/pq                 |
| Rust API   | Rust, axum, sqlx                        |
| Database   | PostgreSQL 17                           |
| Container  | Podman / Docker                         |
| Orchestration | Kubernetes (k3s via Lima on macOS)   |

---

## Project Structure

```
.
├── services/
│   ├── go-api/       # Go REST API
│   ├── rust-api/     # Rust REST API
│   └── frontend/     # React app
├── infra/
│   └── k8s/          # Kubernetes manifests + deployment guide
└── docker-compose.yml
```

---

## Run locally

```bash
podman-compose up --build
# frontend → http://localhost:3000
# go-api   → http://localhost:8080
# rust-api → http://localhost:8081
```

---

## Deploy on Kubernetes

See [`infra/k8s/README.md`](infra/k8s/README.md).