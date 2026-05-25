# Search — searchApi.go

**File:** `Go_Refined_Code/handlers/searchApi.go`

Handles the `GET /api/search` endpoint. Queries the `pages` table using MySQL FULLTEXT search and returns ranked results as JSON. Also tracks every search query in the `search_queries` table so the scraper can prioritise popular terms.

---

## Data Structures

```go
type SearchResult struct {
    Title       string `json:"title"`
    URL         string
    Content     string `json:"content"`
    Description string `json:"description"`
}

type SearchResponse struct {
    Data []SearchResult `json:"data"`
}

type ErrorResponse struct {
    StatusCode int    `json:"statusCode"`
    Message    string `json:"message"`
}
```

`Description` is a truncated (160-character) snippet of `Content`, used by the frontend to show a preview under each result.

---

## Handler: `SearchAPIHandler(db) http.HandlerFunc`

**Route:** `GET /api/search?q=<query>&language=<lang>`

### Flow

1. Reads the `q` query parameter. Returns `422` if it is empty.
2. Reads the `language` parameter. Defaults to `"en"` if not provided.
3. Runs a FULLTEXT search against the `pages` table:

```sql
SELECT title, content, url
FROM pages
WHERE language = ?
  AND (MATCH(title) AGAINST(? IN NATURAL LANGUAGE MODE) > 0
       OR MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE) > 0)
ORDER BY MATCH(title) AGAINST(? IN NATURAL LANGUAGE MODE) * 3
       + MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE) DESC
LIMIT 20
```

Title matches are weighted 3× higher than content matches, so a page whose title contains the search term ranks above a page that only mentions it in the body.

4. Builds a `Description` field for each result by truncating `Content` to 160 runes.
5. Logs the search event as a structured JSON line:
   ```json
   {"time":"...","level":"INFO","msg":"user_search","query":"example","language":"en","result_count":5}
   ```
6. Increments the `search_queries_total` Prometheus counter for the search term.
7. Upserts the query into the `search_queries` table, incrementing the count by 1:
   ```sql
   INSERT INTO search_queries (query, language, count)
   VALUES (?, ?, 1)
   ON DUPLICATE KEY UPDATE count = count + 1
   ```
8. Returns `200` with a `SearchResponse` JSON object.

---

## FULLTEXT Indexes

The `pages` table has two FULLTEXT indexes (added by migration `00006_fulltext_search.sql`):

| Index | Column |
|---|---|
| `ft_title` | `title` |
| `ft_content` | `content` |

These are used separately in the WHERE clause and ORDER BY expression. MySQL's natural language mode automatically ignores common stop words and scores terms by frequency.

---

## Important Notes

- **Maximum 20 results.** The `LIMIT 20` clause is fixed.
- **Language filter.** Every search is filtered by language. There is no cross-language search.
- **No user tracking.** The log event deliberately omits user ID and IP — only the query, language, and result count are recorded.
- **Parameterised query.** The search term is passed as a prepared statement parameter, not interpolated into the SQL string, which prevents SQL injection.
- **Search query tracking feeds the scraper.** The `search_queries` table populated here is read by the Wikipedia scraper to decide which pages to fetch next (see [Scraper](scraper.md)).
