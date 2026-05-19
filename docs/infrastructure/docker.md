# Docker & Containerisation

**Files:**
- `Go_Refined_Code/Dockerfile` — Go backend image
- `Go_Refined_Code/Dockerfile.nginx` — Nginx frontend image
- `Go_Refined_Code/compose.yml` — Production Docker Compose
- `Go_Refined_Code/compose.dev.yml` — Development Docker Compose

---

## Container Images

### Go Backend — `Dockerfile`

Uses a **multi-stage build** to keep the final image small:

**Stage 1 — Build**
```dockerfile
FROM golang:1.26-alpine AS builder
```
- Copies `go.mod` and `go.sum` and downloads modules first (layer cached unless dependencies change).
- Copies the rest of the source code.
- Builds the binary: `go build -o main .`

**Stage 2 — Runtime**
```dockerfile
FROM alpine:latest
```
- Copies only the compiled binary from the build stage.
- No Go toolchain, no source code, no build cache in the final image.
- Exposes port `8080`.
- Entrypoint: `./main`

The resulting image is a few megabytes rather than hundreds.

---

### Nginx Frontend — `Dockerfile.nginx`

```dockerfile
FROM nginx:alpine
```
- Copies the `static/` directory to `/usr/share/nginx/html`.
- Copies `nginx.conf` to `/etc/nginx/conf.d/default.conf`.
- Exposes port `80`.

All routing decisions (which HTML file to serve, which requests to proxy) are handled by `nginx.conf`. See [Nginx](nginx.md) for details.

---

## Docker Compose — Development (`compose.dev.yml`)

Used for local development. Builds images from source code rather than pulling pre-built images.

**Services:**

| Service | Image | Port |
|---|---|---|
| `mysql` | `mysql:8` | 3307 (mapped from 3306) |
| `go-backend` | Built locally from `Dockerfile` | 8080 |
| `frontend-proxy` | `nginx:alpine` with mounted static files | 8081 |

**Notable differences from production:**
- MySQL port is mapped to `3307` on the host to avoid conflicts with a local MySQL installation.
- Static files are mounted as a volume — changes to HTML/CSS/JS appear immediately without rebuilding the image.
- The Go backend is built from local source on `docker compose up`.

**Starting the dev environment:**
```bash
cd Go_Refined_Code
docker compose -f compose.dev.yml up -d
```

---

## Docker Compose — Production (`compose.yml`)

Used on the Azure VM. Pulls pre-built images from GitHub Container Registry.

**Services:**

| Service | Image | Port |
|---|---|---|
| `mysql` | `mysql:8` | 3306 |
| `go-backend` | `ghcr.io/gorillaerne/go-backend:latest` | 8080 |
| `frontend-proxy` | `ghcr.io/gorillaerne/frontend-proxy:latest` | 8081 |
| `postgres` | `postgres:17` | 5432 |

**Persistent volumes:**
- `go_refined_code_mysql_data` — MySQL data directory (survives container restarts).
- `go_refined_code_postgres_data` — PostgreSQL data directory.

**Health checks:**
- MySQL: `mysqladmin ping` every 10 seconds.
- PostgreSQL: `pg_isready` every 10 seconds.
- The `go-backend` service depends on MySQL being healthy before starting.

**Deploying on the server:**
```bash
docker compose pull
docker compose up -d
```

This is run automatically by the CD pipeline via SSH.

---

## Container Networking

All services in both compose files share a Docker bridge network. They communicate using service names as hostnames:
- `go-backend` connects to MySQL using the hostname `mysql`.
- Nginx proxies API requests to `go-backend:8080`.

No ports other than `8081` (Nginx) need to be exposed to the outside world in production.

---

## Image Tags and Registry

The CI/CD pipeline builds and pushes images to `ghcr.io/gorillaerne/`:

| Tag | When created |
|---|---|
| `sha-<commit>` | Every push to `main` |
| `main` | Every push to `main` |
| `latest` | Every push to `main` |
| `v*.*.*` | When a version tag is pushed |

The production `compose.yml` uses `:latest` so `docker compose pull` always fetches the most recently built image.
