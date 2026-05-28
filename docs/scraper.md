# Manually Triggering the Scraper

To index a Wikipedia article, send a POST request to the scrape endpoint.

## Request

```bash
curl -X POST https://gorillahub.dk/api/scrape \
  -H "X-Scrape-Key: <SCRAPE_KEY>" \
  -H "Content-Type: application/json" \
  -d '{"query": "YOUR TOPIC HERE", "language": "en"}'
```

- `query` — the Wikipedia article title to scrape (e.g. `"Python"`, `"World War II"`)
- `language` — `"en"` for English, `"da"` for Danish
- `X-Scrape-Key` — find this in the `.env` file as `SCRAPE_KEY`

## Expected response

```json
{"message": "scrape queued"}
```

## Verify it worked

Wait ~10 seconds, then search for the topic:

```bash
curl "https://gorillahub.dk/api/search?q=YOUR+TOPIC+HERE"
```

Or just search for it on the site.
