# CI/CD Pipelines

**Directory:** `.github/workflows/`

Five GitHub Actions workflows automate security checks, testing, building, and deploying the application. They are chained together so each stage only runs if the previous one succeeded.

```
continuous-integration.yaml
        ↓
continuous-delivery.yaml  →  release.yaml
        ↓
continuous-deployment.yaml
```

`discord-notifications.yaml` runs alongside all workflows and sends Discord alerts on success or failure.

---

## Continuous Integration — `continuous-integration.yaml`

**Triggers:** Push to `main` or `dev`, and pull requests targeting `main` or `dev`.

### Job dependency graph

```
hashpin
    ↓               ↓
go-lint           js-lint
    ↓                 ↓
go-test           playwright-test
    ↓     ↘
  sast    build
```

`sast` and `build` run in parallel after `go-test` passes. The workflow only succeeds when all terminal jobs (sast, build, playwright-test) pass.

---

#### `hashpin` — HashPin Enforcer

Scans every workflow file in `.github/workflows/` and fails if any `uses:` line references a tag or branch instead of a full 40-character commit SHA. This prevents supply chain attacks where a malicious actor could push a new version to a tag you depend on.

All `uses:` lines in this repository must be pinned to a SHA, e.g.:
```yaml
uses: actions/checkout@34e114876b0b11c390a56381ad16ebd13914f8d5 # v4
```

---

#### `go-lint`

Runs [golangci-lint](https://golangci-lint.run) against the Go source code.

- **Needs:** `hashpin`
- **Working directory:** `./Go_Refined_Code`
- **Config file:** `./linters/.golangci.yml`
- **Enabled linters:** `govet`, `errcheck`, `staticcheck`, `unused`, `revive`, `gosec`, `misspell`, `bodyclose`, `gofmt`, `goimports`

---

#### `js-lint`

Runs [Biome](https://biomejs.dev) against the JavaScript files.

- **Needs:** `hashpin`
- **Target:** `./Go_Refined_Code/static/javaScript/`

---

#### `go-test`

Runs the Go test suite and generates a coverage report.

- **Needs:** `go-lint`
- **Command:** `go test -coverprofile=coverage.out ./...`
- Caches the Go module and build cache by hashing `go.sum`.
- Uploads `coverage.out` as an artifact for the `sast` job to consume.

---

#### `playwright-test`

Runs end-to-end browser tests using Chromium.

- **Needs:** `js-lint`
- **Working directory:** `./Go_Refined_Code/tests/e2e`

---

#### `sast` — SonarCloud

Static Application Security Testing. Scans the Go source code for bugs, vulnerabilities, and code smells.

- **Needs:** `go-test` (downloads the coverage artifact)
- **Config:** `Go_Refined_Code/sonar-project.properties`
- **Secrets required:** `SONAR_TOKEN`
- Only reports results for the `main` branch on the free SonarCloud plan.

---

#### `build`

Compiles the Go application to verify there are no build errors.

- **Needs:** `go-test`
- **Command:** `go build ./...`
- Runs in parallel with `sast` and `playwright-test`.

---

## Continuous Delivery — `continuous-delivery.yaml`

**Triggers:** When Continuous Integration completes successfully on `main`, a version tag matching `v*.*.*`, or a manual workflow dispatch.

Builds Docker images, scans them with OWASP ZAP, and only pushes them to GHCR if the scan passes.

### Job 1: `build`

Builds both Docker images in parallel and exports them as gzip tarballs, uploaded as a workflow artifact for subsequent jobs.

```bash
docker build -t ghcr.io/gorillaerne/go-backend:local ./Go_Refined_Code &
docker build -f Go_Refined_Code/Dockerfile.nginx -t ghcr.io/gorillaerne/frontend-proxy:local ./Go_Refined_Code &
wait
```

### Job 2: `dast` — OWASP ZAP

Dynamic Application Security Testing. Loads the backend image, starts it on port 8080, and runs a ZAP full scan against it.

- **Needs:** `build`
- Waits up to 45 seconds for the app to become ready before scanning.
- Uses `fail_action: true` — if ZAP finds any alerts at WARN level or above, the job fails and the image is **not pushed**.
- Creates GitHub Issues for any findings.

### Job 3: `push`

Loads the images from the artifact, logs in to GHCR, and pushes all tags in parallel.

- **Needs:** `dast` — only runs if the DAST scan passed.
- Each image is pushed with three tags simultaneously:
  - `sha-<short-commit>` — immutable reference to the exact build
  - `<branch-or-tag-name>` — e.g. `main` or `v2026.05.27-abc1234`
  - `latest`

---

## Continuous Deployment — `continuous-deployment.yaml`

**Triggers:** When Continuous Delivery completes successfully, or a manual workflow dispatch.

SSHes into the production server and pulls the new images.

```bash
cd /Gorillaerne/Devops_Gorillas/Go_Refined_Code
git fetch origin main && git reset --hard origin/main
docker compose pull
docker compose up -d
```

**Secrets required:** `SERVER_SSH_KEY`, `SERVER_IP`

---

## Release — `release.yaml`

**Triggers:** When Continuous Delivery completes successfully on `main`.

Auto-generates a GitHub Release with a date+SHA tag (e.g. `v2026.05.27-abc1234`) and release notes from merged pull requests since the last release.

Runs after Continuous Delivery so a release only exists once the image has been built, scanned, and pushed to GHCR.

---

## Discord Notifications — `discord-notifications.yaml`

Watches all workflows on `main` and posts to Discord on success or failure.

**Secrets required:** `CI_FAIL_DISCORD_WEBHOOK`, `PUSH_MAIN_DISCORD_WEBHOOK`

---

## Secrets

| Secret | Used by | Description |
|---|---|---|
| `GITHUB_TOKEN` | Delivery, Deployment | Auto-provided by GitHub Actions |
| `SONAR_TOKEN` | CI (sast) | SonarCloud authentication token |
| `SERVER_SSH_KEY` | Deployment | SSH private key for the production server |
| `SERVER_IP` | Deployment | IP address of the production server |
| `CI_FAIL_DISCORD_WEBHOOK` | Discord | Webhook URL for failure notifications |
| `PUSH_MAIN_DISCORD_WEBHOOK` | Discord | Webhook URL for success notifications |

Application secrets (JWT_SECRET, etc.) are set as environment variables directly on the server, not in GitHub.

---

## Pre-Commit Hook

**File:** `.git/hooks/pre-commit`

Runs `golangci-lint` before every commit and blocks the commit if linting fails. This mirrors the `go-lint` job in CI so developers catch issues locally before pushing.
