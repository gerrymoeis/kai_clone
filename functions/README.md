# Cloudflare Pages Functions

This directory contains **Cloudflare Pages Functions** - serverless functions that run alongside your static Pages deployment.

## What are Pages Functions?

Pages Functions are Workers that:
- Deploy automatically with your Pages site (no separate deployment!)
- Match file paths to URLs (`functions/counter/sync.js` → `/counter/sync`)
- Run on Cloudflare's edge network globally
- Are **much simpler** than standalone Workers (no routing config needed)

## Current Functions

### `/counter/sync` - HTMX Counter Demo

**File**: `functions/counter/sync.js`  
**Handles**: `POST /counter/sync`

This function handles the HTMX counter demo from the homepage. It receives a count from the client and echoes it back, demonstrating server-side interaction without a backend server.

**What it does**:
- Receives form data with `count` parameter
- Parses the count as an integer
- Returns the count back as plain text
- HTMX updates the "Server (HTMX)" display

## How It Works

1. **File Structure** → **URL Mapping**
   ```
   functions/counter/sync.js → /counter/sync
   functions/api/users.js    → /api/users
   functions/api/[id].js     → /api/:id (dynamic segment)
   ```

2. **Export Functions for HTTP Methods**
   ```javascript
   export async function onRequestPost(context) { /* POST handler */ }
   export async function onRequestGet(context) { /* GET handler */ }
   export async function onRequest(context) { /* All methods */ }
   ```

3. **Deploy Automatically**
   ```bash
   gforge deploy pages --project=your-project --run
   ```
   Functions are deployed with your Pages site automatically!

## Context Object

The `context` parameter provides:
- `context.request` - The incoming Request
- `context.env` - Environment variables and bindings
- `context.params` - Dynamic route parameters
- `context.waitUntil()` - Run tasks after response
- `context.next()` - Call next middleware

## Example: Add a New Function

**Create**: `functions/api/hello.js`
```javascript
export async function onRequestGet(context) {
  return new Response('Hello from the edge!', {
    headers: { 'Content-Type': 'text/plain' },
  });
}
```

**Deploy**: `gforge deploy pages --project=your-project --run`

**Test**: `curl https://your-project.pages.dev/api/hello`

## Documentation

- [Pages Functions Docs](https://developers.cloudflare.com/pages/functions/)
- [Context Object](https://developers.cloudflare.com/pages/functions/api-reference/)
- [Routing](https://developers.cloudflare.com/pages/functions/routing/)

## Why Pages Functions > Standalone Workers?

| Feature | Pages Functions | Standalone Workers |
|---------|----------------|-------------------|
| **Deployment** | Automatic with Pages | Separate deployment |
| **Routing** | File-based (simple) | Manual config (complex) |
| **Setup** | Zero config | wrangler.toml required |
| **Use Case** | API endpoints for static site | Complex routing, middleware |

**For Gothic Forge demos**: Pages Functions are perfect! ✅
