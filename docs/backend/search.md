# Search — searchApi.go

**File:** `Go_Refined_Code/handlers/searchApi.go`

Handles the `GET /api/search` endpoint. Queries the `pages` table and returns matching results as JSON.

---

## Data Structures

```go
type SearchResult struct {
    Title   string `json:"title"`
    URL     string
    Content string `json:"content"`
}

type SearchResponse struct {
    Data []SearchResult `json:"data"`
}

type ErrorResponse struct {
    StatusCode int    `json:"statusCode"`
    Message    string `json:"message"`
}
```

---

## Handler: `SearchAPIHandler(db) http.HandlerFunc`

**Route:** `GET /api/search?q=<query>&language=<lang>`

### Flow

1. Reads the `q` query parameter. Returns `422` if it is empty.
2. Reads the `language` parameter. Defaults to `"en"` if not provided.
3. Runs the following SQL query against MySQL:

```sql
SELECT title, content, url
FROM pages
WHERE (title LIKE ? OR content LIKE ?)
  AND language = ?
LIMIT 20
```

The search term is wrapped with `%` wildcards so it matches anywhere in the title or content.

4. Scans the rows into `[]SearchResult`.
5. Logs the search event as a structured JSON line:
   ```json
   {"time":"...","level":"INFO","msg":"user_search","query":"example","language":"en","result_count":5}
   ```
6. Increments the `search_queries_total` Prometheus counter for the search term.
7. Returns `200` with a `SearchResponse` JSON object.

---

## Important Notes

- **Maximum 20 results.** The `LIMIT 20` clause is fixed in the query.
- **Language filter.** Every search is filtered by language. There is no way to search across all languages in a single request.
- **No user tracking.** The log event deliberately omits the user ID and IP address — only the query, language, and result count are recorded.
- **Parameterised query.** The search term is passed as a prepared statement parameter, not interpolated into the SQL string, which prevents SQL injection.
