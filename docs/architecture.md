# System Architecture

## Overview

**¿Who Knows?** is a web search engine that lets users search a database of indexed web pages, register accounts, and log in. The project migrates a legacy Python 2 / Flask application (2009) to a modern Go backend.

The system is containerized with Docker and deployed to Microsoft Azure. Three containers work together to serve the application.

---

## High-Level Architecture

```
User's Browser (port 8081)
        │
        ▼
┌───────────────────┐
│   Nginx Container │  — Serves static HTML/CSS/JS
│   (port 8081)     │  — Proxies /api/* to Go backend
└────────┬──────────┘
         │  /api/*
         ▼
┌───────────────────┐
│  Go Backend API   │  — REST API (port 8080)
│  (port 8080)      │  — Authentication, Search, Profile
└────────┬──────────┘
         │
         ▼
┌───────────────────┐
│   MySQL Database  │  — Stores users and indexed pages
│   (port 3306)     │
└───────────────────┘
```

---

## Container Overview

| Container | Image | Port | Responsibility |
|---|---|---|---|
| `frontend-proxy` | `nginx:alpine` | 8081 | Serves static files, proxies API calls |
| `go-backend` | Custom Go binary on `alpine` | 8080 | All business logic and REST API |
| `mysql` | `mysql:8` | 3306 | Persistent data storage |

---

## Request Flow

**A search request from start to finish:**

1. User types a query on the homepage and clicks Search.
2. The browser sends `GET /api/search?q=example` to Nginx on port 8081.
3. Nginx matches the `/api/` prefix and forwards the request to `go-backend:8080`.
4. The Go backend runs a `LIKE` query against the `pages` table in MySQL.
5. Results are returned as a JSON array and rendered in the browser.

**An authenticated request (e.g. change password):**

1. User submits the change-password form on `/profile`.
2. The browser reads the JWT token from `localStorage` and sends it as a `Bearer` header.
3. Nginx proxies the `POST /api/change-password` request to Go.
4. Go validates the JWT, verifies the current password, hashes the new one, and updates the database.

---

## Technology Stack

| Layer | Technology | Why |
|---|---|---|
| Backend language | Go | Performance, static binaries, strong stdlib |
| HTTP router | Gorilla Mux | Pattern matching, subrouters, middleware |
| Database driver | `go-sql-driver/mysql` | Official MySQL driver |
| Authentication | JWT (HS256) | Stateless, no server-side session storage |
| Password hashing | bcrypt (cost 12) | Industry standard, slow by design |
| Metrics | Prometheus | Standard for containerised applications |
| Logging | `log/slog` (JSON) | Structured, machine-parsable log lines |
| Frontend | Vanilla HTML/CSS/JS | No framework — keeps the frontend simple |
| Reverse proxy | Nginx | Efficient static file serving + API proxy |
| Container runtime | Docker / Docker Compose | Local dev and production parity |
| CI/CD | GitHub Actions | Automated lint, test, build, and deploy |
| Hosting | Microsoft Azure VM | Cloud deployment target |
| Container registry | GitHub Container Registry (GHCR) | Free with GitHub |

---

## Directory Layout

```
Devops_Gorillas/
├── Go_Refined_Code/          # Main application (active codebase)
│   ├── database/             # Database connection logic
│   ├── handlers/             # API route handlers + middleware
│   ├── static/               # Frontend files served by Nginx
│   │   ├── html/             # 5 HTML pages
│   │   ├── javaScript/       # 6 JS modules
│   │   └── style.css         # Global stylesheet
│   ├── tests/e2e/            # Playwright end-to-end tests
│   ├── linters/              # Linter configuration files
│   ├── main.go               # Application entry point
│   ├── compose.yml           # Production Docker Compose
│   ├── compose.dev.yml       # Development Docker Compose
│   ├── Dockerfile            # Go backend container
│   ├── Dockerfile.nginx      # Nginx frontend container
│   └── nginx.conf            # Nginx routing config
│
├── Legacy/                   # Original Python 2 Flask app (reference only)
├── .github/workflows/        # CI and CD GitHub Actions pipelines
├── docs/                     # This documentation
└── README.md
```

---

## Environment Variables

The application is configured entirely through environment variables. In development, these are loaded from a `.env` file. In production, they are set on the Docker host.

| Variable | Required | Description |
|---|---|---|
| `DATABASE_PATH` | Yes | MySQL DSN, e.g. `user:pass@tcp(mysql:3306)/whoknowsdb` |
| `JWT_SECRET` | Yes | Secret key used to sign JWT tokens |
| `MYSQL_ROOT_PASSWORD` | Yes | MySQL root password (used by Docker) |
| `MYSQL_DATABASE` | Yes | MySQL database name |
| `MYSQL_USER` | Yes | MySQL application user |
| `MYSQL_PASSWORD` | Yes | MySQL application user password |
| `SEND_BREACH_EMAILS` | No | Set to `true` to send breach notification emails on startup |
| `RESEND_API_KEY` | No | API key for the Resend email service |
| `RESEND_FROM_EMAIL` | No | Sender address for breach notification emails |

---

## Database Schema

### `users` table

| Column | Type | Notes |
|---|---|---|
| `id` | INT, AUTO_INCREMENT, PK | Unique user identifier |
| `username` | TEXT, UNIQUE, NOT NULL | Login name |
| `email` | TEXT, NOT NULL | Used for breach notifications |
| `password` | TEXT, NOT NULL | bcrypt hash (always starts with `$2`) |

### `pages` table

| Column | Type | Notes |
|---|---|---|
| `title` | TEXT | Page title shown in search results |
| `content` | TEXT | Page body — searched with `LIKE` |
| `url` | TEXT | Original URL |
| `language` | TEXT | Language code, e.g. `en` |
