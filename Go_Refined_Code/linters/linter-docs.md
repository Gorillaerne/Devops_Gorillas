## 🔍 Linters & Pre-Commit Hook

**Date:** 23/02/2026

---

We added two linters and a pre-commit hook to enforce code quality before every commit.

### Go Linter (golangci-lint)

Config: `Go_Refined_Code/linters/.golangci.yml`

```bash
# Install
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.9.0

# Run manually
cd Go_Refined_Code
golangci-lint run --config=linters/.golangci.yml ./...
```

### JavaScript Linter (Biome)

Config: `biome.json`

```bash
# macOS
brew install biome@2.4.4

# Windows
choco install biome --version 2.4.4

# Run manually
biome lint ./Go_Refined_Code/static/javaScript/
```

### Pre-Commit Hook

The hook runs both linters automatically on every commit and blocks the commit if either fails. Hook is located in `.githooks/pre-commit`.

**One-time setup per developer:**
```bash
git config core.hooksPath .githooks
```

### Summary

| Feature | Tool | Status |
| :--- | :--- | :--- |
| **Go Linting** | `golangci-lint` | ✅ Completed |
| **JavaScript Linting** | `Biome` | ✅ Completed |
| **Pre-Commit Hook** | Shell script | ✅ Completed |
| **CI/CD Integration** | GitHub Actions | ✅ Completed |
