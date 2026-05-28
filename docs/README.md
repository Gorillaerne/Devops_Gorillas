# Documentation

This folder contains documentation for every part of the ¿Who Knows? system.

---

## Contents

### System Overview
- [Architecture](architecture.md) — How the three containers fit together, technology stack, database schema, and environment variables.
- [API Reference](api-reference.md) — All REST API endpoints, request formats, and response formats.

### Backend (Go)

| File | Documentation |
|---|---|
| `main.go` | [Application Entry Point](backend/main.md) |
| `handlers/authApi.go` + `registerApi.go` | [Authentication](backend/authentication.md) |
| `handlers/searchApi.go` | [Search](backend/search.md) |
| `handlers/profileApi.go` | [Profile / Change Password](backend/profile.md) |
| `handlers/breachList.go` | [Breach Detection](backend/breach-detection.md) |
| `handlers/emailService.go` | [Email Service](backend/email-service.md) |
| `handlers/loggingMiddleware.go` + `metrics.go` | [Middleware & Metrics](backend/middleware.md) |
| `handlers/weatherApi.go` | [Weather (Placeholder)](backend/weather.md) |
| `functions/packages/crawler/crawl/main.go` | [Crawler & Scraper](backend/scraper.md) |
| `database/sqlite.go` | [Database](backend/database.md) |

### Frontend

| File(s) | Documentation |
|---|---|
| `static/html/*.html` | [Pages](frontend/pages.md) |
| `static/javaScript/*.js` | [JavaScript Modules](frontend/javascript.md) |
| `static/style.css` | [Styling](frontend/styling.md) |

### Infrastructure

| File(s) | Documentation |
|---|---|
| `Dockerfile`, `Dockerfile.nginx`, `compose.yml`, `compose.dev.yml` | [Docker & Containers](infrastructure/docker.md) |
| `nginx.conf` | [Nginx](infrastructure/nginx.md) |
| `.github/workflows/` | [CI/CD Pipelines](infrastructure/cicd.md) |
| `linters/.golangci.yml`, `linters/biome.json`, `.githooks/pre-commit` | [Linting & Pre-Commit Hook](infrastructure/linting.md) |

### Testing
- [Testing](testing.md) — Go unit/integration tests and Playwright E2E tests.

---

## Quick Start

```bash
# Clone the repo
git clone https://github.com/Gorillaerne/Devops_Gorillas.git
cd Devops_Gorillas

# Install the pre-commit hook
cp .githooks/pre-commit .git/hooks/pre-commit

# Start the dev environment
cd Go_Refined_Code
docker compose -f compose.dev.yml up -d
```

The application will be available at **http://localhost:8081**.
