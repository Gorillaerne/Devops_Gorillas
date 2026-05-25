# Scraper — scraper.go

**File:** `Go_Refined_Code/handlers/scraper.go`

A background goroutine that periodically fetches Wikipedia articles for the most searched queries and upserts them into the `pages` table. This keeps the search index populated with content users actually want to find.

---

## How It Works

```
Startup
  └── StartScraper(db, 1h)
        └── runScrape() immediately
        └── ticker fires every hour
              └── runScrape()
                    ├── topQueries()  — top 20 from search_queries, skipping recently scraped
                    └── for each query:
                          ├── fetchWikipedia()  — Wikipedia REST API
                          └── upsertPage()      — INSERT ... ON DUPLICATE KEY UPDATE
```

The scraper fires once immediately on startup, then repeats on the configured interval.

---

## `StartScraper(db *sql.DB, interval time.Duration)`

Called from `main.go` after the database is connected:

```go
apiHandlers.StartScraper(database.DB, 1*time.Hour)
```

Starts the background goroutine. The interval is configurable — in production it is set to 1 hour.

---

## Top Queries

```sql
SELECT sq.query, sq.language
FROM search_queries sq
LEFT JOIN pages p ON p.title = sq.query AND p.language = sq.language
WHERE p.last_updated IS NULL OR p.last_updated < NOW() - INTERVAL 24 HOUR
ORDER BY sq.count DESC
LIMIT 20
```

Only queries whose corresponding page is either missing or older than 24 hours are returned. This prevents the scraper from re-fetching the same pages on every run.

---

## Wikipedia API

Pages are fetched from the Wikipedia REST API summary endpoint:

```
https://{language}.wikipedia.org/api/rest_v1/page/summary/{query}
```

- `language` is `en` or `da`, matching the language stored in `search_queries`.
- A `User-Agent` header is sent to identify the application, as required by Wikipedia's API policy.
- A 10-second timeout is applied to each request.

The response provides the article title, a plain-text extract (the introduction), and the canonical page URL.

If Wikipedia returns 404 (no article found) or the extract is empty, the query is skipped with a warning log and the scraper moves on to the next query.

---

## Upsert into Pages

```sql
INSERT INTO pages (title, url, language, content, last_updated)
VALUES (?, ?, ?, ?, NOW())
ON DUPLICATE KEY UPDATE
    title        = VALUES(title),
    content      = VALUES(content),
    last_updated = NOW()
```

The `url` column is the primary key of the `pages` table. If a page already exists, its title and content are refreshed and `last_updated` is set to now — ensuring stale content is kept up to date.

---

## Logging

Each scrape cycle produces structured log output:

| Event | Level | Description |
|---|---|---|
| Scraper started | INFO | Logged once at startup with the configured interval |
| Page upserted | INFO | Logged for each successfully scraped page |
| Wikipedia fetch failed | WARN | Logged when an article is not found or the request fails |
| Top queries fetch failed | ERROR | Logged if the database query for top queries fails |
| Page upsert failed | ERROR | Logged if writing to the `pages` table fails |

---

## Important Notes

- **No duplicate fetches.** The 24-hour window in `topQueries` ensures pages are only re-scraped once per day at most.
- **Failures are non-fatal.** If one query fails (e.g. no Wikipedia article), the scraper logs a warning and continues with the next query. A single failure never stops the whole run.
- **Content is the Wikipedia intro only.** The scraper uses the summary endpoint, which returns the introductory paragraph rather than the full article. This keeps content concise and avoids very large payloads.
- **Language support.** Queries searched in Danish (`da`) are scraped from `da.wikipedia.org`, and English queries from `en.wikipedia.org`.
