# 🚀 Gothic Forge v3 - Quickstart

**Development → Production in 2 Minutes**

---

## 📋 Prerequisites

- **Go 1.21+** - [Download](https://go.dev/dl/)
- **Git** - [Download](https://git-scm.com/downloads)

---

## ⚡ Quick Start (3 Steps)

### **1. Clone & Install**

```bash
git clone https://github.com/gerrymoeis/gothic_forge.git
cd gothic_forge

# Build CLI
go build -o gforge ./cmd/gforge  # or gforge.exe on Windows

# Bootstrap project
./gforge install
```

**What this does**:
- ✅ Installs Go tools (templ, gotailwindcss, air)
- ✅ Creates .env with secure JWT_SECRET
- ✅ Sets up Tailwind + static assets
- ✅ Creates Dockerfile

---

### **2. Start Development**

```bash
# Start dev server with hot reload
./gforge dev
```

Visit: `http://localhost:8080`

**Commands**:
```bash
./gforge test          # Run tests
./gforge build         # Build production binary
./gforge doctor        # Check system health
```

---

### **3. Deploy to Production**

#### **Option A: Static Site Only** (Fastest - 1 minute)

Perfect for: Landing pages, blogs, documentation

```bash
# Export static site
./gforge export

# Deploy to Cloudflare Pages
./gforge deploy pages --project=my-site --run
```

**What you get**:
- ✅ Global CDN (300+ locations)
- ✅ HTMX interactions via Pages Functions
- ✅ Sub-50ms responses
- ✅ 100% free tier

**Setup**: Only needs [Cloudflare API token](https://dash.cloudflare.com/profile/api-tokens)

---

#### **Option B: Full Stack (Leapcell + Cloudflare Proxy)** (Production - ~5 minutes)

Perfect for: Production apps, SaaS, complex applications

```bash
# One-time: CockroachDB service account key (for auto-provision)
./gforge secrets --set COCKROACH_API_KEY=<your-key>

# Optional: Aiven token for Valkey auto-provision (or paste REDIS_URL interactively)
./gforge secrets --set AIVEN_TOKEN=<your-token>

# Deploy (guided)
./gforge deploy --provider=leapcell --with-valkey
```

**What you get**:
- ✅ Go backend (Leapcell compute)
- ✅ PostgreSQL-compatible database (CockroachDB Serverless)
- ✅ Redis-compatible cache (Aiven Valkey)
- ✅ Cloudflare Proxy CDN (orange‑cloud) in front of Leapcell
- ✅ Global distribution on generous free tiers

---

## 🔑 Get API Keys

### **CockroachDB** (Database) ⭐

1. Sign up: [cockroachlabs.cloud/signup](https://cockroachlabs.cloud/signup)
2. Create service account: [Service Accounts](https://cockroachlabs.cloud/service-accounts)
3. Copy API key (shown once!)
4. Save: `./gforge secrets --set COCKROACH_API_KEY=<key>`

**Free tier**: 5 GB storage, perfect for side projects

### **Aiven** (Redis Cache) ⭐

1. Sign up: [console.aiven.io/signup](https://console.aiven.io/signup)
2. Generate token: [Profile → Tokens](https://console.aiven.io/profile/tokens)
3. Copy token (shown once!)
4. Save: `./gforge secrets --set AIVEN_TOKEN=<token>`

**Free trial**: 30 days, then $10/month

### **Cloudflare** (Static Assets) ⭐

1. Sign up: [dash.cloudflare.com/sign-up](https://dash.cloudflare.com/sign-up)
2. Create API token: [Profile → API Tokens](https://dash.cloudflare.com/profile/api-tokens)
   - Template: "Edit Cloudflare Workers"
   - Permissions: Pages Edit + Workers Edit
3. Install CLI: `npm install -g wrangler`
4. Save: `./gforge secrets --set CLOUDFLARE_API_TOKEN=<token>`

**Free tier**: Unlimited static requests, 100k Workers requests/day

### **Back4app** (Go App Hosting) ⭐

**No API key needed!** Uses GitHub integration (guided during deployment).

1. Sign up: [back4app.com/signup](https://www.back4app.com/signup)
2. Connect GitHub when deploying

**Free tier**: 25k container hours/month

---

## 📚 Deployment Paths

### **Path 1: Cloudflare Pages Only**

Best for: Static sites with light interactivity

```bash
./gforge export
./gforge deploy pages --project=my-project --run
```

**Deploy time**: ~1 minute  
**Cost**: $0/month  
**Includes**: Cloudflare Pages Functions for dynamic endpoints

### **Path 2: Full Stack (Railway)**

Best for: Automated deployment, quick setup

```bash
./gforge deploy --with-valkey --with-pages
```

**Deploy time**: ~5 minutes  
**Provider**: Railway (default)  
**Requires**: RAILWAY_TOKEN

### **Path 3: Full Stack (Back4app)**

Best for: Learning Docker/DevOps, git-based deploys

```bash
./gforge deploy --provider=back4app --with-valkey --with-pages
```

**Deploy time**: ~10 minutes (first time)  
**Provider**: Back4app Containers  
**Requires**: Docker installed  
**Re-deploys**: Just `git push`!

---

## 🎯 Which Path is Right for You?

| Need | Choose | Why |
|------|--------|-----|
| **Landing page, blog, docs** | Pages Only | Fastest, free, no backend needed |
| **SaaS, user auth, database** | Full Stack | Production-ready, all features |
| **Quick automated deploy** | Railway | One command, auto-provisions |
| **Learn DevOps workflows** | Back4app | Educational, git-based deploys |

---

## 🛠️ CLI Commands

```bash
# Development
./gforge dev                    # Start dev server (hot reload)
./gforge test                   # Run tests
./gforge test --coverage        # With coverage
./gforge build                  # Build production binary
./gforge doctor                 # Check system health

# Deployment
./gforge deploy --check         # Validate secrets/config
./gforge deploy --dry-run       # Preview without executing
./gforge deploy pages --run     # Deploy static site
./gforge deploy --with-valkey --with-pages  # Full stack

# Database
./gforge db --migrate           # Run migrations
./gforge db --status            # Check migration status

# Secrets
./gforge secrets --gen-jwt      # Generate JWT secret
./gforge secrets --set KEY=val  # Set secret in .env
```

---

## 📦 What Gets Deployed

### **Static Site (Pages Only)**
```
Cloudflare Pages
├── HTML/CSS/JS → CDN (global)
└── functions/
    └── counter/sync.js → Edge endpoint
```

### **Full Stack**
```
Cloudflare Pages (Static Assets)
    ↓
Back4app/Railway (Go Backend)
    ↓
Aiven Valkey (Redis Cache)
    ↓
CockroachDB (PostgreSQL Database)
```

---

## 🔄 CI/CD with GitHub Actions

Gothic Forge includes automated deployment workflows:

**Auto-deploy on push to `main`**:
```yaml
# .github/workflows/ci.yml already configured!
# Just add secrets to GitHub repo:
Settings → Secrets → New repository secret

Required secrets:
- CLOUDFLARE_API_TOKEN
- CF_ACCOUNT_ID
- CF_PROJECT_NAME
- DATABASE_URL (optional)
- RAILWAY_TOKEN (optional)
```

**Manual deploy trigger**:
Go to Actions → Manual Deploy → Run workflow

---

## 🚨 Troubleshooting

### Docker not found (Back4app only)

```bash
# Install Docker Desktop
# Windows/Mac: https://docs.docker.com/desktop/
# Linux: https://docs.docker.com/engine/install/

# Verify
docker --version
```

### Missing tools

```bash
# Re-run install
./gforge install

# Or use doctor with --fix
./gforge doctor --fix
```

### API key invalid

```bash
# Regenerate and set again
./gforge secrets --set KEY_NAME=new-value

# Check .env file directly
cat .env  # Linux/Mac
type .env # Windows
```

---

## 📖 Next Steps

- **Custom Domain**: Cloudflare Pages → Custom Domains
- **Monitor**: View logs at provider dashboards
- **Scale**: Upgrade plans as you grow
- **Features**: Check `CONTRIBUTING.md` for adding features

---

## 🎓 Learning Resources

- **Full Docs**: `README.md`
- **Architecture**: `ARCHITECTURE.md` (coming soon)
- **Contributing**: `CONTRIBUTING.md`
- **Functions Guide**: `functions/README.md`

---

## 📊 Comparison: Deployment Options

| Feature | Pages Only | Railway | Back4app |
|---------|-----------|---------|----------|
| **Setup Time** | 1 min | 5 min | 10 min |
| **Deploy Command** | `deploy pages` | `deploy` | `deploy --provider=back4app` |
| **Re-deploy** | Re-run command | `railway up` | `git push` |
| **Database** | ❌ | ✅ | ✅ |
| **Backend** | ❌ | ✅ | ✅ |
| **Edge Functions** | ✅ | ❌ | ❌ |
| **Cost (Free Tier)** | $0 | $5 credit | 25k hrs |
| **Best For** | Static sites | Quick deploy | Learning |

---

## ⚡ Summary

Gothic Forge gets you from zero to production in **2 commands**:

```bash
# 1. Install
./gforge install

# 2. Deploy
./gforge deploy pages --run          # Static site (1 min)
# OR
./gforge deploy --with-valkey --with-pages  # Full stack (5 min)
```

**That's it!** You're live with:
- ✅ Global CDN
- ✅ HTTPS by default
- ✅ Hot reload in dev
- ✅ Production-ready architecture
- ✅ All on generous free tiers

**Now build something amazing!** 🚀
