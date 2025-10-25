# Deploy Gothic Forge Demo on Cloudflare Only

**Goal**: Deploy the demo app using only Cloudflare (Pages + Functions) - no backend server needed!

**What You Get**:
- ✅ Static site on Cloudflare Pages (global CDN)
- ✅ HTMX counter working (via Pages Functions)
- ✅ Sub-50ms responses worldwide
- ✅ **100% free** (Cloudflare free tier)
- ✅ **No backend server required!**

---

## 🚀 Quick Start (3 Steps)

### **Step 1: Build Static Site**

```bash
# Build the static export (includes Pages Functions)
gforge export

# This creates dist/ with:
# - dist/index.html (homepage)
# - dist/static/ (CSS, JS, images)
# - dist/functions/counter/sync.js (HTMX endpoint)
```

**What happens**:
- Compiles Templ templates to HTML
- Builds Tailwind CSS
- Copies static assets
- **Copies functions/ directory** (Pages Functions)

### **Step 2: Deploy to Cloudflare Pages**

```bash
# Deploy to Pages (interactive, creates project if needed)
gforge deploy pages --project=gothic-forge-demo --run

# Or just:
gforge deploy pages --run
# (wrangler will prompt for project name)
```

**What happens**:
- Uploads dist/ to Cloudflare
- Deploys Pages Functions automatically
- Creates project if it doesn't exist
- Returns deployment URL

**Expected output**:
```
✨ Success! Uploaded 11 files (5.81 sec)
🌎 Deploying...
✨ Deployment complete!
    https://abc123.gothic-forge-demo.pages.dev
```

### **Step 3: Test the Site**

Visit your deployment URL:
```
https://your-project.pages.dev
```

**Test the HTMX counter**:
1. Click the "+" button (local counter)
2. Wait 5 seconds (auto-sync to "server")
3. **Counter syncs via Pages Function!** ✅

---

## 📊 What's Deployed?

| Component | What | Where |
|-----------|------|-------|
| **Static HTML/CSS/JS** | Homepage, styles, Alpine/HTMX | Cloudflare Pages CDN |
| **Pages Function** | `/counter/sync` endpoint | Cloudflare Edge (global) |
| **No Backend** | ❌ No Go server needed! | N/A |

---

## 🔍 How Pages Functions Work

### **File Structure = URL Mapping**

```
functions/counter/sync.js → /counter/sync
```

When you deploy with `gforge deploy pages`, Cloudflare automatically:
1. Detects `dist/functions/` directory
2. Deploys each JS file as an edge function
3. Maps file path to URL path
4. Runs globally on Cloudflare's edge

### **The Function Code**

**File**: `functions/counter/sync.js`
```javascript
export async function onRequestPost(context) {
  const formData = await context.request.formData();
  const count = formData.get('count') || '0';
  const num = parseInt(count, 10) || 0;
  
  return new Response(num.toString(), {
    status: 200,
    headers: { 'Content-Type': 'text/plain; charset=utf-8' },
  });
}
```

**What it does**:
- Receives POST request from HTMX
- Extracts `count` from form data
- Returns count back (server echo)
- Runs in <10ms on the edge

---

## 🎯 Architecture

### **Current (Before)**
```
User → Cloudflare Pages (static only)
  ↓
  ❌ POST /counter/sync → 405 Method Not Allowed
```

### **Now (With Pages Functions)**
```
User → Cloudflare Pages (static HTML/CSS/JS)
  ↓
  GET / → Pages CDN (instant)
  
  POST /counter/sync → Pages Function (edge)
    ↓
    ✅ Returns count in <10ms
```

**No backend server needed!** Everything runs on Cloudflare's edge.

---

## 💡 Benefits

| Metric | Value |
|--------|-------|
| **Response Time** | <50ms (global CDN) |
| **Function Latency** | <10ms (edge execution) |
| **Backend Load** | 0 requests (no backend!) |
| **Cost** | $0/month (free tier) |
| **Regions** | 300+ (Cloudflare network) |
| **Scalability** | Unlimited (edge compute) |

---

## 🛠️ Development Workflow

### **Local Development**

```bash
# Terminal 1: Run Go backend (for local testing)
gforge dev

# Visit: http://localhost:8080
# HTMX counter hits local Go server
```

### **Deploy to Cloudflare**

```bash
# Build + deploy
gforge export
gforge deploy pages --project=your-project --run

# Visit: https://your-project.pages.dev
# HTMX counter hits Pages Function (edge)
```

### **Iterate**

```bash
# Make changes to:
# - app/templates/*.templ (HTML)
# - app/static/app.js (Alpine/HTMX)
# - functions/counter/sync.js (edge function)

# Rebuild + redeploy
gforge export
gforge deploy pages --run
```

---

## 📝 Adding New Functions

### **Example: Add `/api/hello` Endpoint**

**1. Create file**: `functions/api/hello.js`
```javascript
export async function onRequestGet(context) {
  return new Response(
    JSON.stringify({ message: 'Hello from the edge!' }),
    {
      headers: { 'Content-Type': 'application/json' },
    }
  );
}
```

**2. Deploy**:
```bash
gforge export
gforge deploy pages --run
```

**3. Test**:
```bash
curl https://your-project.pages.dev/api/hello
# {"message":"Hello from the edge!"}
```

**That's it!** File-based routing is automatic.

---

## 🎓 Pages Functions vs Standalone Workers

| Feature | Pages Functions | Standalone Workers |
|---------|----------------|-------------------|
| **Deployment** | Automatic with Pages | Separate `wrangler deploy` |
| **Routing** | File-based (zero config) | Manual routes in `wrangler.toml` |
| **Use Case** | API for static site ✅ | Complex apps, middleware |
| **Setup Time** | 0 minutes | 10+ minutes (config) |

**For demos**: Pages Functions are **perfect** ✅

---

## 🚧 Limitations

### **What Pages Functions CAN'T Do** (vs Full Backend)

| Feature | Pages Functions | Go Backend |
|---------|----------------|------------|
| **Database Queries** | ❌ (no direct DB access) | ✅ Full SQL/ORM |
| **Session Management** | ⚠️ (use KV/Durable Objects) | ✅ Native sessions |
| **File Uploads** | ⚠️ (limited, use R2) | ✅ Full filesystem |
| **Long-Running Tasks** | ❌ (30s timeout) | ✅ Background jobs |
| **Complex Business Logic** | ⚠️ (JavaScript only) | ✅ Full Go power |

### **What Pages Functions CAN Do** (Perfectly)

- ✅ API endpoints (JSON responses)
- ✅ Form handling (HTMX, POST data)
- ✅ Redirects and proxying
- ✅ Edge caching
- ✅ Simple data transformation
- ✅ Authentication (JWT validation)
- ✅ Rate limiting (via KV)

**For the demo**: Pages Functions are **sufficient** ✅

---

## 🔮 Future: Full Stack

When you need a full backend:

```bash
# Deploy full stack (Pages + Backend + DB + Cache)
gforge deploy --provider=back4app --with-valkey --with-pages

# Architecture:
# - Cloudflare Pages: Static assets
# - Cloudflare Workers: Edge caching (optional)
# - Back4app: Go backend server
# - Aiven Valkey: Redis cache
# - CockroachDB: PostgreSQL database
```

But for now, **Cloudflare-only deployment works great!** 🎉

---

## ✅ Summary

**What you deployed**:
- ✅ Static site on Cloudflare Pages
- ✅ Pages Function handling `/counter/sync`
- ✅ HTMX counter working globally
- ✅ **No backend server required**
- ✅ **100% free tier**

**Commands**:
```bash
gforge export                                 # Build dist/
gforge deploy pages --project=demo --run     # Deploy to Pages
```

**Result**: Fully working demo app on Cloudflare edge! ⚡

---

## 🎯 Next Steps

1. ✅ **Test your deployment** - Visit the URL
2. ✅ **Add custom domain** (optional) - Cloudflare Dashboard
3. ✅ **Add more functions** - Create files in `functions/`
4. 🔮 **Deploy full stack** (when you need DB/sessions)

Ready to deploy? Run:
```bash
gforge export && gforge deploy pages --run
```

🚀 **Deploy in under 2 minutes!**
