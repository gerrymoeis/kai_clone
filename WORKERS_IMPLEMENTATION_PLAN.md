# Cloudflare Workers Implementation Plan

**Date**: October 24, 2025  
**Status**: 🎯 Design Phase - Ready for Implementation  
**Target**: Gothic Forge v3.1 (Optional Feature)

---

## 🎯 Executive Summary

**Goal**: Add optional Cloudflare Workers layer for edge caching without breaking Gothic Forge's simplicity philosophy.

**Approach**: Annotation-based code generation
- Go handlers annotated with `//gforge:cache` directives
- CLI generates tiny Worker JS files (~30 LOC each)
- Deploy via `gforge deploy --with-workers` (opt-in)
- Zero changes to existing workflows

**Benefits**:
- ⚡ 10-50ms edge responses (vs 150-300ms backend)
- 📉 80-90% traffic offload to edge
- 💰 Massive cost savings at scale
- 🌍 Global distribution automatic

**Philosophy Alignment**: ✅ APPROVED
- Single Go codebase (JS is artifact)
- Opt-in (doesn't break simplicity)
- Teaching through annotations
- Batteries-included tooling

---

## 📐 Architecture

### **Current (v3.0)**
```
User Request
    ↓
Cloudflare Pages (static) ← Global CDN
    ↓
Back4app (Go app) ← Single region (150-300ms)
    ↓
Valkey (Redis) ← Cache layer
    ↓
CockroachDB ← Database
```

**Problems**:
- ❌ All dynamic requests hit backend (expensive at scale)
- ❌ No geographic distribution for dynamic content
- ❌ Backend must handle cache misses + DB queries

### **Proposed (v3.1 with --with-workers)**
```
User Request
    ↓
    ├─→ Static → Cloudflare Pages (CDN)
    │
    └─→ Dynamic → Cloudflare Worker (Edge, global)
            ├─ Cache hit? → Return (10ms) ✅
            │
            └─ Cache miss? → Back4app (Go app)
                    ↓
                Valkey (Redis)
                    ↓
                CockroachDB
```

**Benefits**:
- ✅ 80-90% requests served from edge (10ms)
- ✅ Backend only handles cache misses
- ✅ Global distribution automatic
- ✅ Cost scales linearly with misses, not total traffic

---

## 🔧 Implementation Details

### **Phase 1: Annotation Syntax** (Gothic Forge DSL)

**In Go Handlers**:
```go
package routes

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

// PostFragment renders a single post fragment (HTMX partial)
//
//gforge:cache ttl=60s key=post:{id}:v{version} methods=GET
func PostFragment(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    // Your existing Go logic
    post := db.GetPost(id)
    templ.Render(w, r, PostFragmentTemplate(post))
}

// PostList renders post list (frequently accessed, long TTL)
//
//gforge:cache ttl=5m key=posts:list:page{page} methods=GET vary=Cookie
func PostList(w http.ResponseWriter, r *http.Request) {
    // Your existing Go logic
    posts := db.GetPosts()
    templ.Render(w, r, PostListTemplate(posts))
}

// PostCreate handles form submission (NO caching, write operation)
func PostCreate(w http.ResponseWriter, r *http.Request) {
    // No annotation = Worker passes through (no caching)
    // Your existing Go logic
}
```

**Annotation Format**:
```
//gforge:cache [parameters]

Parameters:
  ttl=<duration>     - Cache TTL (60s, 5m, 1h, 24h)
  key=<pattern>      - Cache key pattern (supports {param} placeholders)
  methods=<list>     - HTTP methods to cache (default: GET)
  vary=<header>      - Vary cache by header (Cookie, Authorization, etc.)
  strategy=<type>    - Cache strategy (default: cache-first)
                       Options: cache-first, stale-while-revalidate, network-first
```

### **Phase 2: Code Generation** (`gforge gen-workers`)

**What it does**:
1. Parses Go files for `//gforge:cache` annotations
2. Extracts route patterns from chi/gorilla router
3. Generates Worker JS for each cached route
4. Creates `wrangler.toml` config
5. Outputs to `workers/` directory

**Generated Worker Example**:
```javascript
// workers/post_fragment.js (GENERATED - DO NOT EDIT)
// Source: app/routes/posts.go:PostFragment
// Cache: ttl=60s key=post:{id}:v{version} methods=GET

addEventListener("fetch", event => {
  event.respondWith(handleRequest(event.request));
});

async function handleRequest(request) {
  const url = new URL(request.url);
  
  // Only cache GET requests
  if (request.method !== 'GET') {
    return fetch(request);
  }
  
  // Extract route params
  const match = url.pathname.match(/^\/posts\/([^/]+)\/fragment$/);
  if (!match) return fetch(request);
  
  const id = match[1];
  
  // Build cache key with version
  const version = await getVersion('post', id); // From KV
  const cacheKey = `post:${id}:v${version}`;
  
  // Try cache first
  const cache = caches.default;
  let response = await cache.match(cacheKey);
  
  if (response) {
    // Cache hit - add header for debugging
    response = new Response(response.body, response);
    response.headers.set('X-Cache', 'HIT');
    response.headers.set('X-Cache-Key', cacheKey);
    return response;
  }
  
  // Cache miss - fetch from backend
  const backendUrl = `${BACKEND_URL}${url.pathname}${url.search}`;
  response = await fetch(backendUrl, {
    method: request.method,
    headers: request.headers,
    body: request.body
  });
  
  // Cache successful responses
  if (response.ok) {
    const cloned = response.clone();
    cloned.headers.set('Cache-Control', 'public, max-age=60');
    await cache.put(cacheKey, cloned);
    
    response = new Response(response.body, response);
    response.headers.set('X-Cache', 'MISS');
    response.headers.set('X-Cache-Key', cacheKey);
  }
  
  return response;
}

// Fetch version from KV (for cache invalidation)
async function getVersion(entity, id) {
  const key = `version:${entity}:${id}`;
  const version = await VERSIONS.get(key);
  return version || '1';
}
```

**Generated wrangler.toml**:
```toml
# workers/wrangler.toml (GENERATED)
name = "gothic-forge-workers"
main = "workers/index.js"
compatibility_date = "2025-01-01"

# Environment variables
[vars]
BACKEND_URL = "https://myapp.b4a.app"

# KV Namespace for version tracking
kv_namespaces = [
  { binding = "VERSIONS", id = "your-kv-namespace-id" }
]

# Routes (automatically mapped from annotations)
routes = [
  { pattern = "yourdomain.com/posts/*/fragment", zone_name = "yourdomain.com" },
  { pattern = "yourdomain.com/posts", zone_name = "yourdomain.com" }
]
```

### **Phase 3: Cache Invalidation** (Backend)

**When data changes, increment version**:

```go
// app/routes/posts.go

//gforge:cache ttl=60s key=post:{id}:v{version} methods=GET
func PostFragment(w http.ResponseWriter, r *http.Request) {
    // ... render post ...
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    // 1. Update database
    err := db.UpdatePost(id, newData)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    
    // 2. Invalidate edge cache (increment version)
    workers.IncrementVersion("post", id)
    
    // 3. Respond
    w.WriteHeader(200)
}
```

**Workers SDK** (auto-generated):
```go
// internal/workers/invalidation.go (GENERATED)

package workers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "os"
)

// IncrementVersion increments version in Cloudflare KV
// This invalidates all cached responses with old version
func IncrementVersion(entity, id string) error {
    if os.Getenv("GFORGE_WORKERS_ENABLED") != "true" {
        return nil // No-op if Workers not enabled
    }
    
    kvNamespace := os.Getenv("CLOUDFLARE_KV_NAMESPACE_ID")
    apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
    accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
    
    key := "version:" + entity + ":" + id
    
    // Fetch current version
    currentVersion := getKVValue(kvNamespace, key, apiToken, accountID)
    newVersion := currentVersion + 1
    
    // Write new version
    return putKVValue(kvNamespace, key, newVersion, apiToken, accountID)
}

// Helper functions
func getKVValue(namespace, key, token, account string) int {
    // HTTP GET to Cloudflare KV API
    // Returns current version or 0 if not exists
}

func putKVValue(namespace, key string, value int, token, account string) error {
    // HTTP PUT to Cloudflare KV API
}
```

### **Phase 4: CLI Commands**

**New Commands**:
```bash
# Generate Workers from annotations
gforge gen-workers

# Preview generated Workers (dry-run)
gforge gen-workers --dry-run

# Deploy with Workers
gforge deploy --with-workers

# Deploy Workers only (after backend is deployed)
gforge deploy workers

# Test Workers locally (uses miniflare/wrangler dev)
gforge workers dev

# Check Workers status
gforge workers status

# Purge all edge caches (emergency)
gforge workers purge

# Show cache analytics
gforge workers stats
```

**Deploy Flow**:
```bash
# Option A: Full stack + Workers
gforge deploy --provider=back4app --with-workers --with-valkey --with-pages

# Runs:
# 1. gen-workers (parse annotations, generate JS)
# 2. Build Go app
# 3. Deploy to Back4app
# 4. Deploy Workers to Cloudflare
# 5. Deploy Pages to Cloudflare
# 6. Provision Valkey
# 7. Provision CockroachDB

# Option B: Just update Workers (after code change)
gforge deploy workers

# Option C: Deploy without Workers (default)
gforge deploy --provider=back4app --with-valkey --with-pages
```

---

## 📊 Performance Metrics (Expected)

### **Before Workers** (Current v3.0)
| Request Type | Latency | Cached | Traffic |
|--------------|---------|--------|---------|
| Static (HTML/CSS/JS) | 50ms | ✅ Pages CDN | 30% |
| Dynamic (HTMX) - First | 200ms | ❌ | 70% |
| Dynamic (HTMX) - Cached | 20ms | ✅ Valkey | 70% |

**Backend Load**: 70% of total traffic

### **After Workers** (v3.1 with --with-workers)
| Request Type | Latency | Cached | Traffic |
|--------------|---------|--------|---------|
| Static (HTML/CSS/JS) | 50ms | ✅ Pages CDN | 30% |
| Dynamic (HTMX) - Edge Hit | **10ms** | ✅ **Workers** | **60%** ✨ |
| Dynamic (HTMX) - Edge Miss | 220ms | ❌ | 10% |
| Dynamic (HTMX) - Valkey Hit | 20ms | ✅ Valkey | 10% |

**Backend Load**: 10-20% of total traffic (80-90% offload!) 🎉

### **Cost Impact** (10k monthly users, 100k requests)
| Metric | Without Workers | With Workers | Savings |
|--------|----------------|--------------|---------|
| **Backend Requests** | 70k | 10k | **86% reduction** |
| **Backend Cost** | $10/mo | $2/mo | **$8/mo saved** |
| **Workers Cost** | $0 | $0 | Free tier (100k/day) |
| **Response Time (p50)** | 150ms | **30ms** | **5x faster** |
| **Response Time (p95)** | 300ms | **50ms** | **6x faster** |

---

## 🎓 Developer Experience

### **Example: Adding Cache to New Route**

**Before** (no caching):
```go
func NewFeature(w http.ResponseWriter, r *http.Request) {
    data := fetchData()
    templ.Render(w, r, NewFeatureTemplate(data))
}
```

**After** (with edge caching):
```go
//gforge:cache ttl=5m key=feature:new methods=GET
func NewFeature(w http.ResponseWriter, r *http.Request) {
    data := fetchData()
    templ.Render(w, r, NewFeatureTemplate(data))
}
```

**Deploy**:
```bash
gforge deploy --with-workers
```

**Result**:
- ✅ Worker automatically generated
- ✅ Deployed to edge
- ✅ 5-minute cache enabled
- ✅ Backend traffic reduced

**That's it!** One annotation, one deploy.

### **Local Development**

```bash
# Terminal 1: Run Go backend
gforge dev

# Terminal 2: Run Workers locally (proxies to localhost:8080)
gforge workers dev

# Browser: localhost:8787 (Workers) → localhost:8080 (Go)
```

**Testing**:
- ✅ Test with Workers locally
- ✅ Test without Workers (direct Go)
- ✅ Compare performance
- ✅ Debug both layers

---

## 🚧 Implementation Phases

### **Phase 1: Core Infrastructure** (2-3 weeks)
**Goal**: Basic Worker generation + deployment

**Tasks**:
- [ ] Parser for `//gforge:cache` annotations
- [ ] Route extraction from chi router
- [ ] Worker JS template generator
- [ ] `gforge gen-workers` command
- [ ] `wrangler.toml` generator
- [ ] Basic deployment (`gforge deploy workers`)

**Deliverables**:
- ✅ Annotation syntax defined
- ✅ Code generation working
- ✅ Manual deployment possible
- ✅ Example Workers generated

### **Phase 2: Cache Invalidation** (1-2 weeks)
**Goal**: Backend can invalidate edge caches

**Tasks**:
- [ ] KV namespace management
- [ ] Version tracking system
- [ ] `workers.IncrementVersion()` SDK
- [ ] Cloudflare API integration
- [ ] `gforge workers purge` command

**Deliverables**:
- ✅ Backend can invalidate caches
- ✅ Version-based cache keys
- ✅ Emergency purge working

### **Phase 3: Local Development** (1 week)
**Goal**: Test Workers locally

**Tasks**:
- [ ] Miniflare integration
- [ ] `gforge workers dev` command
- [ ] Proxy to local Go backend
- [ ] Hot reload on code changes

**Deliverables**:
- ✅ Local Workers testing
- ✅ Fast iteration cycle
- ✅ Debug experience

### **Phase 4: Monitoring & Analytics** (1 week)
**Goal**: Visibility into Workers performance

**Tasks**:
- [ ] `X-Cache: HIT/MISS` headers
- [ ] `gforge workers stats` command
- [ ] Cache hit ratio reporting
- [ ] Error tracking

**Deliverables**:
- ✅ Cache hit ratio visible
- ✅ Performance metrics
- ✅ Error debugging

### **Phase 5: Documentation** (1 week)
**Goal**: Users can adopt Workers easily

**Tasks**:
- [ ] Annotation reference guide
- [ ] Caching strategies guide
- [ ] Troubleshooting guide
- [ ] Migration guide

**Deliverables**:
- ✅ Complete documentation
- ✅ Example apps
- ✅ Best practices

**Total Timeline**: 6-8 weeks for full implementation

---

## 🎯 Rollout Strategy

### **v3.1.0-alpha** (Week 6)
**Scope**: Basic generation + deployment
- ✅ `gforge gen-workers` works
- ✅ Manual deployment possible
- ⚠️ No invalidation yet
- ⚠️ Experimental, opt-in

**Users**: Early adopters, testing

### **v3.1.0-beta** (Week 8)
**Scope**: Full feature set
- ✅ Cache invalidation working
- ✅ Local development support
- ✅ Monitoring + analytics
- ⚠️ Breaking changes possible

**Users**: Beta testers, dogfooding

### **v3.1.0-stable** (Week 10)
**Scope**: Production-ready
- ✅ All features complete
- ✅ Documentation complete
- ✅ Tested at scale
- ✅ Stable API

**Users**: General availability

### **v3.2.0+** (Future)
**Advanced Features**:
- [ ] Per-user edge caching
- [ ] GraphQL support
- [ ] WebSocket proxying
- [ ] A/B testing at edge
- [ ] Bot detection
- [ ] DDoS protection

---

## ⚖️ Trade-offs Analysis

### **Pros** ✅

| Benefit | Impact |
|---------|--------|
| **80-90% traffic offload** | Backend scales effortlessly |
| **10-50ms edge responses** | 5-10x faster UX |
| **Global distribution** | Low latency worldwide |
| **Cost savings** | 80%+ backend cost reduction |
| **Single codebase** | Go remains source of truth |
| **Opt-in** | Doesn't break existing apps |
| **Teaching value** | Learn edge caching patterns |

### **Cons** ❌

| Trade-off | Mitigation |
|-----------|------------|
| **Added complexity** | Optional, well-documented |
| **JS in toolchain** | Generated, not hand-written |
| **Debugging two layers** | Local Workers dev support |
| **Cache invalidation complexity** | Version-based strategy (simple) |
| **Cloudflare lock-in** | Workers are tiny, portable |

### **Verdict**: **Benefits >> Costs** (for most apps)

**When to use Workers**:
- ✅ Public-facing apps (marketing sites, blogs, docs)
- ✅ High traffic (10k+ requests/day)
- ✅ Read-heavy workloads (80%+ GETs)
- ✅ Global audience

**When to skip Workers**:
- ❌ Internal tools (low traffic)
- ❌ Highly personalized (per-user content)
- ❌ Mostly writes (no cache benefit)
- ❌ Local/single-region users

---

## 🔐 Security Considerations

### **1. Authentication**
**Problem**: Workers can't see encrypted session cookies  
**Solution**: 
- Public endpoints: Cache freely
- Authenticated endpoints: `vary=Cookie` or bypass Workers

### **2. Secrets**
**Problem**: Backend secrets shouldn't be in Workers  
**Solution**: Workers only cache responses, never access secrets

### **3. Rate Limiting**
**Problem**: Edge caching bypasses rate limits  
**Solution**: Implement rate limiting in Workers or backend

### **4. DDoS**
**Problem**: Cache poisoning attacks  
**Solution**: Validate cache keys, use signed URLs for sensitive data

---

## 📚 Example Apps

### **Example 1: Blog** (Perfect for Workers)
```go
//gforge:cache ttl=1h key=post:{slug} methods=GET
func BlogPost(w http.ResponseWriter, r *http.Request) {
    slug := chi.URLParam(r, "slug")
    post := db.GetPost(slug)
    templ.Render(w, r, BlogPostTemplate(post))
}

//gforge:cache ttl=5m key=blog:index methods=GET
func BlogIndex(w http.ResponseWriter, r *http.Request) {
    posts := db.GetPosts()
    templ.Render(w, r, BlogIndexTemplate(posts))
}
```

**Result**: 95% traffic served from edge, sub-10ms responses

### **Example 2: E-commerce** (Selective caching)
```go
//gforge:cache ttl=10m key=product:{id} methods=GET
func ProductDetail(w http.ResponseWriter, r *http.Request) {
    // Product pages cached (10 min)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
    // No annotation = no caching (writes always go to backend)
}

//gforge:cache ttl=1m key=cart:{session} methods=GET vary=Cookie
func ViewCart(w http.ResponseWriter, r *http.Request) {
    // Per-session caching (1 min)
}
```

**Result**: Product pages fast, cart always fresh

---

## 🎓 Learning Path for Users

### **Beginner** (Week 1)
1. Deploy without Workers (understand baseline)
2. Add one annotation to a read-heavy route
3. Deploy with `--with-workers`
4. Check `X-Cache` headers in browser
5. Compare response times

### **Intermediate** (Week 2-3)
1. Add cache invalidation to write operations
2. Test cache purging works
3. Monitor cache hit ratio
4. Optimize TTL values

### **Advanced** (Week 4+)
1. Per-user edge caching
2. Vary by headers
3. Stale-while-revalidate patterns
4. Custom invalidation strategies

---

## ✅ Go/No-Go Decision

### **Recommendation**: **GO** ✅

**Reasoning**:
1. ✅ Aligns with philosophy (opt-in, generated, teaching)
2. ✅ Massive performance benefit (5-10x faster)
3. ✅ Cost savings (80%+ backend offload)
4. ✅ Competitive advantage (most Go frameworks don't have this)
5. ✅ Implementation is tractable (6-8 weeks)
6. ✅ Reversible (can disable anytime)

**Conditions**:
- ⚠️ Must remain optional (flag-gated)
- ⚠️ Must have excellent docs (tutorials + examples)
- ⚠️ Must work locally (wrangler dev)
- ⚠️ Must be stable (no breaking changes post-v3.1.0)

---

## 📝 Next Steps

### **Immediate** (This Week)
1. [ ] Create spike: basic Worker generation
2. [ ] Test Cloudflare API integration
3. [ ] Validate wrangler deployment works
4. [ ] Prototype annotation parser

### **Short-term** (Next 2 Weeks)
1. [ ] Implement Phase 1 (core infrastructure)
2. [ ] Create example app with Workers
3. [ ] Write initial documentation
4. [ ] Test with real backend

### **Medium-term** (Next 4 Weeks)
1. [ ] Implement Phase 2 (cache invalidation)
2. [ ] Implement Phase 3 (local dev)
3. [ ] Implement Phase 4 (monitoring)
4. [ ] Beta release

### **Long-term** (Next 8 Weeks)
1. [ ] Production release (v3.1.0)
2. [ ] Gather user feedback
3. [ ] Iterate on edge cases
4. [ ] Plan v3.2.0 features

---

## 🎉 Summary

**The annotation-based Workers approach is a perfect fit for Gothic Forge.**

**Why it works**:
- ✅ Preserves "one Go codebase" philosophy
- ✅ Adds massive performance benefit
- ✅ Opt-in (doesn't break simplicity)
- ✅ Teaching value (annotations document caching)
- ✅ Implementation is realistic (6-8 weeks)

**What makes it different**:
- Not a split codebase (Go + manual JS)
- Not complex WASM compilation
- Just simple annotations + codegen
- Generated Workers are tiny and readable

**Bottom line**: This is the right way to add edge caching to Gothic Forge. Let's do it! 🚀

**Target**: v3.1.0 (6-8 weeks)
