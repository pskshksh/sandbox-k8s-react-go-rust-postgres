## Architecture

```
  Browser
     │
     ▼
 [Ingress]          ← single entry point, port 80
     │
 [frontend]         ← nginx serves React, proxies /go-api/* and /rust-api/*
     │         │
 [go-api]   [rust-api]    ← ClusterIP, not exposed externally
     │         │
    [postgres]             ← StatefulSet with persistent volume
```

All resources live in a single namespace: `sandbox`.

---

## K8s structure

```
infra/k8s/
├── namespace.yaml
├── secret.yaml
├── network-policy.yaml
├── ingress.yaml
├── postgres/
│   ├── statefulset.yaml
│   └── service.yaml
├── go-api/
│   ├── deployment.yaml
│   └── service.yaml
├── rust-api/
│   ├── deployment.yaml
│   └── service.yaml
└── frontend/
    ├── deployment.yaml
    └── service.yaml
```