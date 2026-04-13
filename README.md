# ¿Who Knows?

A modernized web search engine — a migration of the original 2009 "Who Knows?" Flask application from Python 2 to Go.

The application lets users search a database of web pages, register an account, and log in. It is deployed on Microsoft Azure and served through Nginx.

## Architecture

Three-tier containerized stack:

```
Browser → Nginx (port 8081) → Go API (port 8080) → MySQL
```

- **Nginx** — serves static HTML/CSS/JS and reverse-proxies `/api/*` to the Go backend
- **Go backend** — REST API built with Gorilla Mux, JWT auth, bcrypt password hashing
- **MySQL** — stores `users` and `pages` tables

The active codebase lives in `Go_Refined_Code/`. The original Python app is preserved read-only in `Legacy/src/`.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Go 1.25+](https://go.dev/dl/) (for local development without Docker)
- [golangci-lint](https://golangci-lint.run/usage/install/) (required by the pre-commit hook)
- [Biome](https://biomejs.dev/guides/getting-started/) (required by the pre-commit hook)

## Getting started

### 1. Clone the repository

```bash
git clone https://github.com/Gorillaerne/Devops_Gorillas.git
cd Devops_Gorillas
```

### 2. Set up the pre-commit hook

The hook runs `golangci-lint` and Biome before every commit. It must be active or your commits will not be linted.

```bash
cp .githooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

### 3. Start the dev stack

```bash
cd Go_Refined_Code
docker compose -f compose.dev.yml up -d
```

This starts:
- MySQL on host port `3307`
- Go backend on port `8080`
- Nginx on port `8081`

Open [http://localhost:8081](http://localhost:8081) in your browser.

### Environment variables

Secrets are loaded from a `.env` file in `Go_Refined_Code/` for local development. Create one based on the following:

```env
DATABASE_PATH=user:password@tcp(mysql:3306)/whoknows
JWT_SECRET=your-secret-here
```

In production, these are set as environment variables directly.

## Running tests

### Go tests (unit + integration)

```bash
cd Go_Refined_Code
go test ./...
```

Uses an in-memory SQLite database — no running MySQL required.

### Frontend tests (Playwright)

```bash
cd Go_Refined_Code/tests/e2e
npm install
npx playwright install --with-deps chromium
npx playwright test
```

### Linting

```bash
# Go
cd Go_Refined_Code
golangci-lint run --config=./linters/.golangci.yml

# JavaScript
biome lint ./Go_Refined_Code/static/javaScript/
```

## Contributing

### Workflow

1. Pick an issue from the [project board](https://github.com/orgs/Gorillaerne/projects/1) and move it to **In Progress**.
2. Pull latest `main` and create a feature branch:

   ```bash
   git checkout main && git pull origin main
   git checkout -b your-branch-name
   ```

3. Make sure the pre-commit hook is active (see [Setup](#2-set-up-the-pre-commit-hook) above).
4. Write your changes.
5. Before pushing, verify all of the following pass:
   - `go build ./...` compiles
   - `go test ./...` passes
   - `npx playwright test` passes (from `Go_Refined_Code/tests/e2e/`)
   - All API endpoints respond correctly against the dev stack
   - All frontend pages load and function in the browser
6. Push your branch and open a pull request to `main`, filling out the PR template.

### API endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/search?q=<query>&language=<lang>` | Search pages |
| `POST` | `/api/register` | Register a new user |
| `POST` | `/api/login` | Log in, returns JWT |
| `GET` | `/api/weather` | Weather stub |
| `GET` | `/api/logout` | Logout stub |

## Deployment

The application is deployed to an Azure VM (`51.120.83.21`) via GitHub Actions on every push to `main`. The CD pipeline builds Docker images, pushes them to GitHub Container Registry (`ghcr.io/gorillaerne/`), and deploys over SSH.

See `.github/workflows/` for the full CI/CD pipeline.

## Legacy application

The original Python 2 Flask app is preserved in `Legacy/src/` for reference only. See [`Legacy/src/README.md`](Legacy/src/README.md) for its setup instructions.
