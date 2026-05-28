# Linting & Pre-Commit Hook

**Files:**
- `Go_Refined_Code/linters/.golangci.yml` — Go linter configuration
- `Go_Refined_Code/linters/biome.json` — JavaScript linter configuration
- `.githooks/pre-commit` — Pre-commit hook script

Two linters enforce code quality: **golangci-lint** for Go and **Biome** for JavaScript. Both run automatically on every commit via a git pre-commit hook, and again in CI on every push and pull request.

---

## Go Linter — golangci-lint

**Config:** `Go_Refined_Code/linters/.golangci.yml`

golangci-lint runs multiple linters in a single pass and is significantly faster than running them individually.

### Enabled Linters

| Linter | What It Catches |
|---|---|
| `govet` | Suspicious code constructs (e.g. wrong `Printf` format strings) |
| `errcheck` | Unhandled error return values |
| `staticcheck` | Bugs, performance issues, and simplification opportunities |
| `unused` | Unexported functions, types, and variables that are never used |
| `revive` | General Go style and best practice violations |
| `gosec` | Potential security issues (SQL injection, weak crypto, etc.) |
| `misspell` | Spelling mistakes in comments and strings |
| `bodyclose` | HTTP response bodies that are never closed (causes connection leaks) |

### Enabled Formatters

| Formatter | What It Enforces |
|---|---|
| `gofmt` | Standard Go formatting (indentation, spacing) |
| `goimports` | Correct import grouping and removal of unused imports |

### Other Settings

- **Timeout:** 5 minutes (allows for large codebases without false timeouts)
- **Runs on tests:** Yes (`tests: true`) — test files are linted too
- **Issue limits removed:** `max-issues-per-linter: 0` and `max-same-issues: 0` ensure all issues are reported, not just the first few

### Installing golangci-lint

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.9.0
```

### Running Manually

```bash
cd Go_Refined_Code
golangci-lint run --config=linters/.golangci.yml ./...
```

---

## JavaScript Linter — Biome

**Config:** `Go_Refined_Code/linters/biome.json`

Biome is a fast, all-in-one JavaScript toolchain that handles both linting and formatting.

### Scope

- **Included:** `Go_Refined_Code/static/**/*.js`
- **Excluded:** HTML files

### Enabled Rules

In addition to Biome's built-in `recommended` ruleset, these rules are explicitly enforced:

| Rule | Category | Severity | What It Catches |
|---|---|---|---|
| `noUnusedVariables` | correctness | error | Variables declared but never used |
| `noUnreachable` | correctness | error | Code after a `return`/`throw` that can never run |
| `useIsNan` | correctness | error | Using `=== NaN` instead of `isNaN()` |
| `noDuplicateCase` | suspicious | error | Duplicate `case` labels in a `switch` |
| `noEmptyBlockStatements` | suspicious | error | Empty `{}` blocks with no code inside |
| `noUselessConstructor` | complexity | error | Class constructors that do nothing |

### Installing Biome

**macOS:**
```bash
brew install biome@2.4.4
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://github.com/biomejs/biome/releases/download/@biomejs/biome@2.4.4/biome-win32-x64.exe" -OutFile "biome.exe"
New-Item -ItemType Directory -Path "$HOME\bin" -Force
Move-Item .\biome.exe "$HOME\bin\biome.exe"
[Environment]::SetEnvironmentVariable("Path", "$env:Path;$HOME\bin", "User")
```

### Running Manually

```bash
biome lint ./Go_Refined_Code/static/javaScript/
```

---

## Pre-Commit Hook

**File:** `.githooks/pre-commit`

A shell script that runs both linters before every `git commit`. If either linter fails, the commit is blocked with an error message.

### What the Hook Does

```
git commit
    └── pre-commit hook runs
        ├── cd Go_Refined_Code && golangci-lint run --config=linters/.golangci.yml ./...
        │   └── FAIL → print error message, exit 1 (commit blocked)
        │   └── PASS → continue
        └── biome lint ./Go_Refined_Code/static/javaScript/
            └── FAIL → print "Biome failed. Commit blocked.", exit 1
            └── PASS → "All checks passed. Proceeding with commit."
```

The hook is cross-platform: it calls `biome` on macOS/Linux and `biome.exe` on Windows.

> **Note:** The hook lints the entire codebase on each commit, not just the changed files. This guarantees that no pre-existing lint failures are silently carried forward.

### One-Time Setup (per developer)

```bash
git config core.hooksPath .githooks
```

This tells git to use the `.githooks/` directory for hook scripts. The hook file is already present in the repository — you only need to run this command once after cloning.

Alternatively, copy the hook manually:
```bash
cp .githooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit   # macOS/Linux only
```

---

## Linting in CI

Both linters also run in the CI pipeline on every push and pull request, even if a developer skips the pre-commit hook. See [CI/CD Pipelines](cicd.md) for details on the `go-lint` and `js-lint` jobs.

### Summary

| Check | Tool | Runs locally | Runs in CI |
|---|---|---|---|
| Go linting | golangci-lint | Pre-commit hook | `go-lint` job |
| JavaScript linting | Biome | Pre-commit hook | `js-lint` job |
| Go formatting | gofmt / goimports | Pre-commit hook | `go-lint` job |
