# CI/CD Pipelines

**Directory:** `.github/workflows/`

Four GitHub Actions workflows automate testing, building, and deploying the application.

---

## CI Pipeline — `CI.yaml`

**Triggers:** Push to `main` or `dev`, and pull requests targeting `main` or `dev`.

This pipeline ensures that code is always linted, tested, and buildable before it can be merged.

### Jobs (run in parallel where possible)

```
go-lint ──────┐
              ├── go-test ──┐
js-lint ──────┘             ├── build
              ├── playwright─┘
              └── (js-lint must pass before playwright)
```

---

#### `go-lint`

Runs [golangci-lint](https://golangci-lint.run) against the Go source code.

- **Working directory:** `./Go_Refined_Code`
- **Config file:** `./linters/.golangci.yml`
- **Enabled linters:** `govet`, `errcheck`, `staticcheck`, `unused`, `revive`, `gosec`, `misspell`, `bodyclose`, `gofmt`, `goimports`

The linter catches code style issues, unused variables, unhandled errors, and potential security problems before they reach the codebase.

---

#### `js-lint`

Runs [Biome](https://biomejs.dev) against the JavaScript files.

- **Target:** `./Go_Refined_Code/static/javaScript/`
- Biome enforces consistent code style and catches common JavaScript mistakes.

---

#### `go-test`

Runs the Go unit and integration test suite.

- **Needs:** `go-lint` to pass first.
- **Command:** `go test ./...`
- **Setup:** Caches the Go module cache and build cache by hashing `go.sum` so dependencies are not re-downloaded on every run.
- Tests use in-memory SQLite so no database container is required.

---

#### `playwright-test`

Runs end-to-end browser tests.

- **Needs:** `js-lint` to pass first.
- **Working directory:** `./Go_Refined_Code/tests/e2e`
- **Steps:**
  1. Install Node.js 20 and npm dependencies (`npm ci`).
  2. Install Chromium browser and system dependencies (`npx playwright install --with-deps chromium`).
  3. Run tests (`npx playwright test`).
- Tests spin up a local static file server and drive a real Chromium browser.

---

#### `build`

Compiles the Go application to catch any build errors that tests might miss.

- **Needs:** All four previous jobs to pass.
- **Command:** `go build ./...`

---

## CD Pipeline — `CD.yaml`

**Triggers:** Push to `main`, a version tag matching `v*.*.*`, or a manual workflow dispatch.

Builds Docker images, pushes them to the GitHub Container Registry, and deploys them to the Azure server.

### Job 1: `build-and-push`

1. **Logs in to GHCR** using the `GITHUB_TOKEN` secret (automatically available in GitHub Actions).
2. **Builds and pushes the Go backend image** tagged as:
   - `ghcr.io/gorillaerne/go-backend:sha-<commit>`
   - `ghcr.io/gorillaerne/go-backend:main`
   - `ghcr.io/gorillaerne/go-backend:latest`
3. **Builds and pushes the Nginx frontend image** with the same tag pattern.

Uses `docker/build-push-action` with layer caching enabled (cached by GitHub Actions cache) to speed up repeated builds.

### Job 2: `deploy`

**Needs:** `build-and-push` to complete successfully.

Connects to the Azure VM over SSH and runs the deployment:

```bash
cd /path/to/Devops_Gorillas
git pull origin main
docker compose pull
docker compose up -d
```

This pulls the newly built images from GHCR and restarts the containers. Existing containers are replaced with zero-downtime (Docker Compose starts the new container before stopping the old one by default).

**Secrets required:**
- `SSH_PRIVATE_KEY` — private key to authenticate with the Azure VM.
- `AZURE_VM_IP` — IP address of the server (`51.120.83.21`).

---

## Other Workflows

### `release.yaml`

Triggers on version tags (`v*.*.*`). Creates a GitHub Release automatically, allowing the team to publish official versioned releases from the repository.

### `discord.yaml`

Sends a notification to a Discord channel when a workflow completes. Useful for the team to get notified of CI failures or successful deployments without watching GitHub directly.

---

## Secrets and Environment Variables

The following secrets must be set in the GitHub repository settings:

| Secret | Used by | Description |
|---|---|---|
| `GITHUB_TOKEN` | CD | Auto-provided by GitHub Actions for GHCR auth |
| `SSH_PRIVATE_KEY` | CD | SSH key for Azure VM access |
| `AZURE_VM_IP` | CD | Azure VM hostname or IP |

Application secrets (JWT_SECRET, database credentials) are set as environment variables directly on the Azure VM, not in GitHub.

---

## Pre-Commit Hook

**File:** `.githooks/pre-commit`

A local git hook that runs before every commit. It runs `golangci-lint` and `biome lint` so that developers catch issues before pushing rather than waiting for CI.

To install it:
```bash
cp .githooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```
