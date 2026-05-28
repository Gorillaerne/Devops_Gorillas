# Crawler & Scraper — DigitalOcean Function

**File:** `functions/packages/crawler/crawl/main.go`
**Config:** `functions/project.yml`

A DigitalOcean serverless function that runs on a cron schedule, crawls Wikipedia starting from the most searched queries, follows links to discover related articles, and upserts all discovered content into the `pages` table. It is invoked by DigitalOcean's scheduler — the application server plays no part in triggering or running it.

---

## How It Works

The system has two distinct layers inside the function:

- **Crawler** — manages a BFS queue, discovers new pages by following Wikipedia links
- **Scraper** — fetches article content from the Wikipedia REST API and writes it to the database

```
DigitalOcean scheduler (cron: "0 * * * *")
  └── Main()                          — function entry point
        ├── connects to MySQL via DATABASE_PATH secret
        └── runCrawl()
              ├── topQueries()          — top 10 seeds from search_queries
              └── BFS queue loop:
                    ├── fetchWikipedia()          — scraper: REST API summary
                    ├── upsertPage()              — scraper: INSERT … ON DUPLICATE KEY UPDATE
                    └── fetchWikipediaLinks()     — crawler: MediaWiki links API
                          └── enqueue linked titles at depth+1
```

The function is triggered by DigitalOcean's built-in scheduler once per hour. It does not run inside the application server.

---

## Constants

| Constant | Value | Purpose |
|---|---|---|
| `crawlerTopN` | 10 | Number of top queries used as seeds |
| `crawlerTimeout` | 10s | HTTP timeout per request |
| `crawlerMaxDepth` | 1 | How many hops to follow from a seed (0 = seeds only) |
| `crawlerMaxLinks` | 5 | Maximum links to follow per page |

Values are deliberately sized to keep a full run within the 300s DigitalOcean Functions limit (max ~60 pages × 2 requests = ~120 HTTP calls). Tune in `functions/packages/crawler/crawl/main.go`.

---

## Scheduling

The function is triggered by DigitalOcean's scheduler via `functions/project.yml`:

```yaml
triggers:
  - name: hourly-crawl
    sourceType: scheduler
    sourceDetails:
      cron: "0 * * * *"
    function: crawler/crawl
```

To change the schedule, update the `cron` field and redeploy with `doctl serverless deploy functions/`.

---

## Crawler: Seed Queries

```sql
SELECT sq.query, sq.language
FROM search_queries sq
LEFT JOIN pages p ON p.title = sq.query AND p.language = sq.language
WHERE p.last_updated IS NULL OR p.last_updated < NOW() - INTERVAL 24 HOUR
ORDER BY sq.count DESC
LIMIT 20
```

Only queries whose corresponding page is either missing or older than 24 hours are returned as seeds. This avoids re-fetching recently updated pages.

---

## Crawler: Link Discovery

After a page is scraped, the crawler fetches its linked articles using the MediaWiki API:

```
https://{language}.wikipedia.org/w/api.php?action=query&titles={title}&prop=links&pllimit=10&plnamespace=0&format=json
```

- `plnamespace=0` restricts results to article pages only (no talk pages, categories, etc.).
- Up to `crawlerMaxLinks` (10) linked titles are enqueued at `depth + 1`.
- A `visited` map prevents the same page from being fetched twice within a single crawl run.
- Pages at `crawlerMaxDepth` are scraped but their links are not followed.

---

## Scraper: Wikipedia Summary API

Pages are fetched from the Wikipedia REST API summary endpoint:

```
https://{language}.wikipedia.org/api/rest_v1/page/summary/{title}
```

- `language` is `en` or `da`, matching the language stored in `search_queries`.
- A `User-Agent` header is sent, as required by Wikipedia's API policy.
- A 10-second timeout is applied to each request.

The response provides the article title, a plain-text extract (the introduction), and the canonical page URL.

If Wikipedia returns 404 or the extract is empty, the item is skipped with a warning log and the crawler continues.

---

## Scraper: Upsert into Pages

```sql
INSERT INTO pages (title, url, language, content, last_updated)
VALUES (?, ?, ?, ?, NOW())
ON DUPLICATE KEY UPDATE
    title        = VALUES(title),
    content      = VALUES(content),
    last_updated = NOW()
```

If a page already exists, its title and content are refreshed and `last_updated` is set to now.

---

## Logging

| Event | Level | Description |
|---|---|---|
| Crawler started | INFO | Logged once at startup with the configured interval |
| Page upserted | INFO | Logged for each successfully scraped page, including depth |
| Wikipedia fetch failed | WARN | Logged when an article is not found or the request fails |
| Links fetch failed | WARN | Logged when the MediaWiki links API call fails |
| Top queries fetch failed | ERROR | Logged if the database query for top queries fails |
| Page upsert failed | ERROR | Logged if writing to the `pages` table fails |

---

## Important Notes

- **Crawler + scraper separation.** `fetchWikipedia` and `upsertPage` are the scraper layer — they are called for every discovered page regardless of how it was found. `fetchWikipediaLinks` is the crawler layer — it only runs when depth allows.
- **No duplicate fetches within a run.** The `visited` map ensures each `language:title` pair is processed at most once per crawl cycle.
- **No duplicate fetches across runs.** The 24-hour window in `topQueries` ensures seed pages are only re-seeded once per day at most. Linked pages discovered mid-crawl are not subject to this filter — they are controlled solely by the `visited` map for the current run.
- **Failures are non-fatal.** A failed fetch or upsert logs a warning/error and the crawler moves on to the next item in the queue.
- **Content is the Wikipedia intro only.** The scraper uses the summary endpoint, which returns the introductory paragraph rather than the full article.
- **Language support.** Seeds searched in Danish (`da`) crawl `da.wikipedia.org`, and English seeds crawl `en.wikipedia.org`. Linked pages inherit the language of the page they were discovered from.
