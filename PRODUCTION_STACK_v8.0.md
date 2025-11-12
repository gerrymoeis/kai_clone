# Gothic Forge v8.0 Production Stack

This document describes the recommended production architecture for Gothic Forge v8.0 and how to configure Cloudflare as a proxy CDN in front of your Leapcell monolith.

## Stack Summary

- Compute: Leapcell (monolith, Go Fiber + HTMX + Tailwind + Templ)
- CDN/Proxy: Cloudflare (orange-cloud proxied DNS)
- Static-only optional: Cloudflare Pages (via `gforge deploy pages`)
- Database: CockroachDB Serverless (Basic)
- Cache: Aiven Valkey (Redis-compatible)

## Goals

- Zero-cost/free-tier friendly
- Global performance with Cloudflare edge caching
- Simple DX: monolith-first, optional Pages for static-only sections

## Cloudflare Proxy Configuration (Free)

1) Add your domain to Cloudflare and change nameservers (one-time).
2) In DNS, point your app hostname (e.g., app.example.com) to your Leapcell URL/IP and ensure the orange cloud is ON (Proxied).
3) Create Cache Rules (or Page Rules if Cache Rules not available). Keep it to 3 simple rules:

- Rule 1: Cache Everything for static assets
  - If URL path matches: /static/*
  - Edge TTL: 1 year
  - Respect origin Cache-Control (we send `public, max-age=31536000, immutable`).

- Rule 2: Cache Everything for public HTML pages
  - If URL path matches: / and any public sections (e.g., /docs/*)
  - Edge TTL: 60–300 seconds
  - Respect origin Cache-Control (we send `public, s-maxage=60, stale-while-revalidate=300` by default).

- Rule 3: Bypass for dynamic/auth/api
  - If URL path matches: /api/*, /auth/*, /dashboard/* (add more as needed)
  - Cache: Bypass

Notes
- Do NOT cache authenticated pages. The server sets `private, no-store` for HTMX/with-session responses.
- The server adds `Vary: Cookie, Accept` where needed to avoid cache poisoning.

## Server Cache Headers (already implemented)

- Static: `Cache-Control: public, max-age=31536000, immutable`
- Public HTML (no session): `Cache-Control: public, s-maxage=60, stale-while-revalidate=300`
- HTMX or session-present: `Cache-Control: private, no-store`

Configure via env:
- `DISABLE_HTML_CACHE=0|1`
- `CACHE_PUBLIC_TTL=60` (seconds)
- `CACHE_SWREVAL_TTL=300` (seconds)

## Optional: Cloudflare Pages (static-only)

Use `gforge deploy pages` to deploy a static export to Cloudflare Pages.
This is not required for the monolith path and should only be used for purely static sections.

## CockroachDB Serverless (Basic)

- Free monthly resources: 50 million RUs and 10 GiB storage
- Docs/pricing: https://www.cockroachlabs.com/pricing/

## Aiven Valkey

- Free/Developer plan suitable for small apps; single-node service with backups.
- Effective dataset size is constrained by the service `maxmemory` (fraction of RAM); refer to Aiven docs.
- Docs: Free plan overview https://aiven.io/docs/platform/concepts/free-plan
- Valkey memory usage: https://aiven.io/docs/products/valkey/concepts/memory-usage

## Verification

Use curl to verify headers at the origin (or via Cloudflare once proxied):

```
curl -I https://yourdomain.com/
curl -I https://yourdomain.com/static/styles/output.css
curl -I -H "HX-Request: true" https://yourdomain.com/some-fragment
```

Expectations
- HTML (public): `Cache-Control: public, s-maxage=60, stale-while-revalidate=300`
- Static: `Cache-Control: public, max-age=31536000, immutable`
- HTMX/session or auth: `Cache-Control: private, no-store`

## Operational Notes

- To quickly disable HTML caching in incidents, set `DISABLE_HTML_CACHE=1` and redeploy.
- Keep Cloudflare rules minimal to avoid unexpected edge behavior on Free tier.
- Monitor origin request rate via Leapcell and cache hit ratio via Cloudflare analytics.
