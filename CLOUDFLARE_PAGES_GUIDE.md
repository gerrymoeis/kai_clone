# Cloudflare Pages Deployment Guide

**Deploy Static Sites with Edge Functions in 1 Minute**

---

## 🎯 When to Use Cloudflare Pages Only

Perfect for:
- ✅ Landing pages, marketing sites
- ✅ Blogs, documentation sites
- ✅ Portfolios, showcase sites
- ✅ Static sites with **light interactivity** (forms, counters, simple APIs)

**Not suitable for**:
- ❌ Complex authentication systems
- ❌ Database-heavy applications
- ❌ File uploads (use R2 separately)
- ❌ Long-running background jobs

---

## ⚡ Quick Deploy (3 Steps)

### **1. Build Static Site**

```bash
./gforge export
```

This creates `dist/` with:
- HTML/CSS/JS (static assets)
- `functions/` directory (edge functions)

### **2. Get Cloudflare API Token**

1. Sign up: [dash.cloudflare.com/sign-up](https://dash.cloudflare.com/sign-up)
2. Go to: [Profile → API Tokens](https://dash.cloudflare.com/profile/api-tokens)
3. Create token with permissions:
   - **Account → Cloudflare Pages → Edit**
   - **Zone → Workers Scripts → Edit**
4. Save token:
   ```bash
   ./gforge secrets --set CLOUDFLARE_API_TOKEN=<your-token>
   ./gforge secrets --set CF_ACCOUNT_ID=<your-account-id>
   ./gforge secrets --set CF_PROJECT_NAME=my-project
   ```

### **3. Deploy**

```bash
./gforge deploy pages --project=my-project --run
```

**Done!** Your site is live at `https://my-project.pages.dev`

---

## 🔧 Cloudflare Pages Functions

**What are they?**
- Serverless functions that run on Cloudflare's edge
- Deploy automatically with your static site
- File-based routing (zero config!)
- Run globally in 300+ locations

**Example**: The HTMX counter demo uses a Pages Function!

### **Current Functions**

#### `/counter/sync` - HTMX Counter Endpoint

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
- Receives POST from HTMX
- Echoes count back (server interaction demo)
- Runs in <10ms on the edge

---

## 📝 Adding New Functions

### **Simple API Endpoint**

**Create**: `functions/api/hello.js`

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

**Deploy**: `./gforge export && ./gforge deploy pages --run`

**Test**: `curl https://your-project.pages.dev/api/hello`

### **Form Handler**

**Create**: `functions/contact/submit.js`

```javascript
export async function onRequestPost(context) {
  try {
    const formData = await context.request.formData();
    const email = formData.get('email');
    const message = formData.get('message');
    
    // Send to email service (e.g., SendGrid, Resend)
    // Or save to KV storage
    
    return new Response('Thank you! Message received.', {
      status: 200,
      headers: { 'Content-Type': 'text/plain' },
    });
  } catch (error) {
    return new Response('Error processing form', {
      status: 400,
    });
  }
}
```

### **Dynamic Route with Parameters**

**Create**: `functions/posts/[id].js`

```javascript
export async function onRequestGet(context) {
  const postId = context.params.id;
  
  // Fetch from KV, external API, or return static data
  const post = {
    id: postId,
    title: `Post ${postId}`,
    content: 'This is dynamically generated!',
  };
  
  return new Response(
    JSON.stringify(post),
    {
      headers: { 'Content-Type': 'application/json' },
    }
  );
}
```

**URL**: `/posts/123` → `context.params.id === "123"`

---

## 🔄 Current Workflow (Manual)

### **For Backend Logic**

**Gothic Forge Philosophy**: We value developer choice and simplicity.

**Current Process** (Manual - V3):
1. Write backend logic in Go (`app/routes/`)
2. Test locally with `./gforge dev`
3. Manually create equivalent JS function in `functions/`
4. Deploy with `./gforge deploy pages`

**Example**:

**Go Code** (`app/routes/api.go`):
```go
func GetUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    user := db.GetUser(id)
    json.NewEncoder(w).Encode(user)
}
```

**Equivalent Pages Function** (`functions/api/users/[id].js`):
```javascript
export async function onRequestGet(context) {
  const id = context.params.id;
  
  // Fetch from external API or KV
  const user = await fetchUserFromAPI(id);
  
  return new Response(JSON.stringify(user), {
    headers: { 'Content-Type': 'application/json' },
  });
}
```

---

## 🚀 Future: Auto-Compilation (V4)

**Planned Feature**: Automatic Go → JavaScript compilation

**Vision**:
```go
// app/routes/api.go

//gforge:edge ttl=60s
func GetUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    user := db.GetUser(id)
    json.NewEncoder(w).Encode(user)
}
```

**Result**: `./gforge gen-edge` would automatically generate:
```javascript
// functions/api/users/[id].js (GENERATED)
export async function onRequestGet(context) {
  // Auto-generated from Go code
}
```

**Why not now?**
- Complex to implement (requires annotation parser, AST transformation)
- Current manual approach is simple and transparent
- V3 focuses on core features and stability
- V4 will add advanced edge compilation

**Current approach benefits**:
- ✅ Simple and understandable
- ✅ Full control over edge functions
- ✅ Easy to debug
- ✅ No magic or black boxes

---

## 📊 Pages Functions Capabilities

### **✅ What You CAN Do**

| Feature | Supported | Notes |
|---------|-----------|-------|
| **API Endpoints** | ✅ | JSON responses, REST APIs |
| **Form Handling** | ✅ | POST data, file uploads (small files) |
| **Redirects** | ✅ | 301/302 redirects |
| **Edge Caching** | ✅ | Cache-Control headers |
| **CORS** | ✅ | Custom headers |
| **Authentication** | ⚠️ | JWT validation only (no session storage) |
| **KV Storage** | ✅ | Requires binding setup |
| **External APIs** | ✅ | Fetch to external services |

### **❌ What You CANNOT Do**

| Feature | Limitation | Alternative |
|---------|-----------|------------|
| **Database Queries** | ❌ Direct SQL | Use external API or full backend |
| **Session Storage** | ❌ Stateful sessions | Use KV or external Redis |
| **File System** | ❌ Local files | Use R2 or external storage |
| **Long Processing** | ❌ >30s timeout | Use full backend or Queue |
| **WebSockets** | ❌ | Use Durable Objects (paid) |

---

## 🎯 Decision Tree: Which Deployment?

```
Need database queries? 
├─ YES → Full Stack (Back4app/Railway + CockroachDB)
└─ NO → Next question

Need user authentication/sessions?
├─ YES → Full Stack (Backend + Redis)
└─ NO → Next question

Need file uploads?
├─ YES → Full Stack + R2
└─ NO → Next question

Just static site + light interactivity?
└─ YES → Cloudflare Pages Only ✅
```

---

## 📦 Example Projects

### **1. Blog with Comments** (Pages Only)

**Structure**:
```
functions/
├── api/
│   ├── posts/
│   │   └── [id].js          # Get post by ID
│   └── comments/
│       └── submit.js         # Submit comment
└── _middleware.js            # CORS, auth
```

**Comments stored in**: KV, external API, or Airtable

### **2. Contact Form** (Pages Only)

**Structure**:
```
functions/
└── contact/
    └── submit.js             # Send email via SendGrid
```

### **3. Analytics Dashboard** (Pages Only)

**Structure**:
```
functions/
└── api/
    ├── events.js             # Track events
    └── stats.js              # Get statistics
```

**Data stored in**: KV or Analytics Engine

---

## 🔧 Testing Functions Locally

**Option 1: With Full Backend**

```bash
# Run Go backend (handles all routes)
./gforge dev

# Test at: http://localhost:8080/counter/sync
```

**Option 2: With Wrangler (Cloudflare CLI)**

```bash
# Install wrangler
npm install -g wrangler

# Test Pages Functions locally
wrangler pages dev dist/
```

---

## 🚀 CI/CD Integration

Pages Functions deploy automatically with GitHub Actions:

**`.github/workflows/ci.yml`** (already configured):
```yaml
- name: Export static site (SSG)
  run: go run ./cmd/gforge export -o dist

- name: Deploy to Cloudflare Pages
  uses: cloudflare/pages-action@v1
  with:
    apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
    accountId: ${{ secrets.CF_ACCOUNT_ID }}
    projectName: ${{ secrets.CF_PROJECT_NAME }}
    directory: dist  # Includes functions/!
```

**Functions are deployed automatically!** No extra steps needed.

---

## 📖 Resources

- **Cloudflare Pages Docs**: [developers.cloudflare.com/pages](https://developers.cloudflare.com/pages/)
- **Functions API**: [developers.cloudflare.com/pages/functions](https://developers.cloudflare.com/pages/functions/)
- **Gothic Forge Functions**: `functions/README.md`
- **Examples**: `functions/counter/sync.js`

---

## 💡 Best Practices

1. **Keep functions small** - <10 KB per function
2. **Cache aggressively** - Use Cache-Control headers
3. **Minimize dependencies** - Functions run on edge, size matters
4. **Use KV for data** - Bind KV namespaces for storage
5. **Handle errors** - Always return proper status codes
6. **Test locally** - Use `wrangler pages dev` before deploying
7. **Monitor** - Check Cloudflare dashboard for errors

---

## 🎯 Summary

**Cloudflare Pages deployment**:
- ✅ **Perfect for**: Static sites with light interactivity
- ✅ **Deploy time**: 1-2 minutes
- ✅ **Cost**: $0 (free tier)
- ✅ **Performance**: <50ms globally
- ✅ **Functions**: File-based routing, zero config

**Commands**:
```bash
./gforge export                          # Build site + functions
./gforge deploy pages --project=NAME --run   # Deploy
```

**When to upgrade to full stack**:
- Need PostgreSQL database
- Need session storage
- Need file uploads (large)
- Need background jobs

**Gothic Forge Philosophy**: 
- Start simple (Pages Only)
- Upgrade when needed (Full Stack)
- Always transparent, never magic
- Your choice, your control

---

**Ready to deploy?** 🚀

```bash
./gforge export && ./gforge deploy pages --run
```
