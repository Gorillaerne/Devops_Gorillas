# Testing

The project has two levels of testing: Go unit/integration tests and Playwright end-to-end browser tests.

---

## Go Tests

**Directory:** `Go_Refined_Code/handlers/`

Test files live alongside the handler files they test. They use the standard Go `testing` package and an **in-memory SQLite database** so tests run without a MySQL server.

### Test Files

| File | What It Tests |
|---|---|
| `auth_test.go` | Password hashing, registration, login |
| `search_test.go` | Search query logic and response format |
| `profile_test.go` | JWT-authenticated password change |
| `breach_test.go` | Breach detection (`isBreached`) |
| `email_test.go` | Email sending via Resend API (mocked server) |

### Test Setup — `TestMain`

Each test file relies on `TestMain` to:
1. Set `bcryptCost = bcrypt.MinCost` — bcrypt at cost 4 instead of 12. This makes password hashing ~100x faster in tests without affecting correctness.
2. Create an in-memory SQLite database with the same schema as the MySQL production database.
3. Pass the database handle to the handler functions under test.

### Running the Tests

```bash
cd Go_Refined_Code
go test ./...
```

To see verbose output:
```bash
go test ./... -v
```

To run tests for a specific file:
```bash
go test ./handlers/ -run TestSearch
```

### What Is Tested

- **auth_test.go** — Verifies that registration creates a user with a bcrypt hash, that login succeeds with correct credentials, and fails with wrong credentials or a non-existent user.
- **search_test.go** — Verifies that search returns results matching the query, respects the language filter, and returns an empty list when there are no matches.
- **profile_test.go** — Verifies that a valid JWT token allows a password change, that an invalid token is rejected, and that the wrong current password is rejected.
- **breach_test.go** — Verifies that known breached credentials are flagged and that unknown credentials are not.
- **email_test.go** — Starts a local HTTP server to mock the Resend API, then verifies that `SendBreachNotification` sends the expected HTTP request.

---

## Playwright End-to-End Tests

**Directory:** `Go_Refined_Code/tests/e2e/`

E2E tests drive a real Chromium browser through the application's UI to verify that the pages work correctly from the user's perspective.

### Configuration — `playwright.config.js`

- **Base URL:** `http://localhost:4321`
- **Browser:** Chromium only
- **Web server:** The config starts `npx serve` before tests run, which serves the `static/` directory on port 4321.

### Test Files

| File | What It Tests |
|---|---|
| `frontend.spec.js` | Search page interactions, navigation links |
| `profile.spec.js` | Profile page — password change form behaviour |

### What Is Tested

- **frontend.spec.js** — Navigates to the homepage, types a search query, clicks the search button, and verifies results appear. Also checks that the Login and Register nav links exist.
- **profile.spec.js** — Navigates to the profile page, verifies the redirect to `/login` when not authenticated, and tests the password change form layout and interactions.

### Running E2E Tests Locally

```bash
cd Go_Refined_Code/tests/e2e

# Install dependencies (first time only)
npm ci
npx playwright install --with-deps chromium

# Run tests
npx playwright test

# Open the Playwright test report
npx playwright show-report
```

### Running in CI

The Playwright tests run in the `playwright-test` job in the CI pipeline (see [CI/CD](infrastructure/cicd.md)). The job installs Node.js, installs dependencies, installs Chromium, and runs `npx playwright test`.

---

## Testing Strategy

| Concern | Covered By |
|---|---|
| Password hashing correctness | Go unit tests |
| SQL query logic | Go integration tests (SQLite) |
| JWT generation and validation | Go unit tests |
| Breach detection logic | Go unit tests |
| Email API integration | Go tests with mocked HTTP server |
| UI rendering and user flows | Playwright E2E tests |
| Code quality and style | golangci-lint, Biome |
| Build success | `go build ./...` in CI |

The Go tests catch regressions in business logic quickly and cheaply. Playwright tests verify the full user experience including the browser-side JavaScript, which Go tests cannot cover.
