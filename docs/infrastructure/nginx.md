# Nginx — nginx.conf

**File:** `Go_Refined_Code/nginx.conf`

Nginx acts as the single entry point for all traffic on port 8081. It serves static files and proxies API requests to the Go backend.

---

## What Nginx Does

1. **Serves static HTML files** for each page URL.
2. **Proxies all `/api/*` requests** to the Go backend container.

---

## Route Configuration

```
GET /              → /usr/share/nginx/html/html/index.html
GET /login         → /usr/share/nginx/html/html/login.html
GET /register      → /usr/share/nginx/html/html/register.html
GET /about         → /usr/share/nginx/html/html/about.html
GET /profile       → /usr/share/nginx/html/html/profile.html
GET /weather       → /usr/share/nginx/html/html/index.html  (fallback)
ANY /api/*         → http://go-backend:8080  (reverse proxy)
```

If a requested file does not exist for a page route, Nginx returns a `404` rather than falling back to an index page.

---

## API Proxy Configuration

```nginx
location /api/ {
    proxy_pass http://go-backend:8080;
    proxy_set_header Host              $host;
    proxy_set_header X-Real-IP         $remote_addr;
    proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
}
```

- `go-backend` resolves to the Go container via Docker's internal DNS.
- The three `proxy_set_header` lines preserve the real client IP and host so the Go backend can log them correctly. Without these, every request would appear to come from Nginx's IP.

---

## Static File Root

All static files (HTML, CSS, JS, images) are copied into the container at `/usr/share/nginx/html` during the Docker build:

```
/usr/share/nginx/html/
├── html/
│   ├── index.html
│   ├── login.html
│   ├── register.html
│   ├── profile.html
│   └── about.html
├── javaScript/
│   └── *.js
├── style.css
└── *.png
```

---

## Why Nginx Instead of Serving Files From Go?

- **Separation of concerns** — Go handles business logic; Nginx handles file serving and routing.
- **Performance** — Nginx is highly optimised for serving static files; there is no overhead from Go's HTTP handler.
- **Flexibility** — Future caching headers, TLS termination, rate limiting, or gzip compression can be added to `nginx.conf` without touching Go code.
