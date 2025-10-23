# 🚀 Gothic Forge v3 - Quickstart Guide

**From Zero to Production in ~15 Minutes**

This guide walks you through cloning Gothic Forge, setting up your development environment, and deploying to production using our **opinionated omakase stack**.

---

## 📋 Prerequisites

- **Go 1.21+** - [Download](https://go.dev/dl/)
- **Git** - [Download](https://git-scm.com/downloads)
- **Node.js 18+** (for deployment CLIs) - [Download](https://nodejs.org/)

**Optional** (Gothic Forge can install these):
- `templ` - Template engine (auto-installed)
- `gotailwindcss` - Tailwind CSS compiler (auto-installed)
- Docker Desktop - Only needed for Back4app Containers ([Download](https://docs.docker.com/desktop/install/windows-install/))

---

## 🎯 Opinionated Stack (Recommended)

Gothic Forge is opinionated about best-in-class services:

| Service | Purpose | Why We Chose It |
|---------|---------|-----------------|
| **CockroachDB Serverless** | Database | PostgreSQL-compatible, global distribution, true serverless, free tier |
| **Aiven Valkey** | Cache/Redis | Redis 7.2+ compatible, managed, high performance, free tier |
| **Cloudflare Pages** | Static Assets | Global CDN, instant deploys, generous free tier |
| **Cloudflare Workers** | Serverless Functions | Edge computing, low latency worldwide |
| **Back4app Containers** | Go App Hosting | Easy container deployment, GitHub integration, free tier |

**Why These?** - Generous free tiers, excellent performance, minimal config, production-ready

---

## 🏁 Step 1: Clone & Bootstrap

### 1.1 Clone the Repository

```bash
# HTTPS
git clone https://github.com/gerrymoeis/gothic_forge.git
cd gothic_forge

# OR SSH
git clone git@github.com:gerrymoeis/gothic_forge.git
cd gothic_forge
```

### 1.2 Check Your Environment

```bash
# Run the doctor to see what you have/need
gforge doctor
```

**Expected Output:**
```
Gothic Forge v3 :: CLI
Doctor
  • Go: go1.24.5 (windows/amd64)
  • git: C:\Program Files\Git\cmd\git.exe
  • templ: (missing)
  • gotailwindcss: (missing)
  • railway: (missing)
  • wrangler: (missing)
  • docker: (missing)
```

### 1.3 Bootstrap Your Project

```bash
# Install all dependencies and tools
gforge install

# This will:
# ✅ Install templ (template engine)
# ✅ Install gotailwindcss (CSS compiler)
# ✅ Generate Tailwind config
# ✅ Create .env from .env.example
# ✅ Generate strong JWT_SECRET
# ✅ Create Dockerfile & .dockerignore
```

**What If Tools Are Missing?**

Gothic Forge attempts to auto-install Go tools. If installation fails, you'll get:

```bash
⚠️  templ: not found
    Install: go install github.com/a-h/templ/cmd/templ@latest
```

Run the suggested command manually.

---

## 🛠️ Step 2: Development

### 2.1 Start Dev Server with Hot Reload

```bash
gforge dev
```

This starts:
- **Templ watch** - Recompiles templates on change
- **Tailwind watch** - Rebuilds CSS on change  
- **Air** - Reloads Go app on change

Visit: **http://localhost:8080**

**Available Routes:**
- `/` - Home page
- `/counter` - HTMX counter demo
- `/healthz` - Health check (liveness)
- `/readyz` - Readiness check (dependencies)

### 2.2 Test Your Build

```bash
# Build production binary
gforge build

# Run the binary
./bin/server

# Or on Windows
.\bin\server.exe
```

### 2.3 Run Tests

```bash
# Run all tests with pretty output
gforge test

# Run with coverage
gforge test --coverage

# Verbose mode
gforge test -v
```

---

## 🔐 Step 3: Configure Secrets

### 3.1 Generate JWT Secret (Required)

```bash
gforge secrets --gen-jwt
```

This creates/updates `.env` with a cryptographically secure JWT_SECRET.

### 3.2 Get API Keys for Deployment

You'll need API keys for each service. Gothic Forge will guide you, but here's a quick reference:

#### **CockroachDB Serverless** (Database) - RECOMMENDED ⭐

1. **Sign up**: [https://cockroachlabs.cloud/signup](https://cockroachlabs.cloud/signup)
2. **Create API key**: 
   - Go to [Account → API Access](https://cockroachlabs.cloud/account/api-access)
   - Click **"Create API Key"**
   - Name: `gothic-forge-deploy`
   - **Copy the key** (shown only once!)
3. **Save to .env**:
   ```bash
   gforge secrets --set COCKROACH_API_KEY=<your-key-here>
   ```

**Free Tier**: 5 GB storage, 1 vCPU, 50M Request Units/month - Perfect for side projects!

#### **Aiven Valkey** (Redis Cache) - RECOMMENDED ⭐

1. **Sign up**: [https://console.aiven.io/signup](https://console.aiven.io/signup)
2. **Create token**: 
   - Go to [Account → Tokens](https://console.aiven.io/account/tokens)
   - Click **"Generate token"**
   - Description: `gothic-forge-deploy`
   - Permissions: **"Read/Write"**
   - **Copy the token**
3. **Save to .env**:
   ```bash
   gforge secrets --set AIVEN_TOKEN=<your-token>
   ```

**Free Trial**: 30 days, then $10/month for Startup plan

#### **Cloudflare Pages** (Static Assets) - RECOMMENDED ⭐

1. **Sign up**: [https://dash.cloudflare.com/sign-up](https://dash.cloudflare.com/sign-up)
2. **Create API token**: 
   - Go to [Profile → API Tokens](https://dash.cloudflare.com/profile/api-tokens)
   - Click **"Create Token"**
   - Template: **"Edit Cloudflare Workers"**
   - Permissions: 
     - Account → Cloudflare Pages → **Edit**
     - Zone → Workers Scripts → **Edit**
   - **Copy the token**
3. **Install Wrangler** (Cloudflare CLI):
   ```bash
   npm install -g wrangler
   ```
4. **Save to .env**:
   ```bash
   gforge secrets --set CLOUDFLARE_API_TOKEN=<your-token>
   gforge secrets --set CF_PROJECT_NAME=gothic-forge-demo
   ```

**Free Tier**: Unlimited static requests, 100k Workers requests/day

#### **Back4app Containers** (Go App) - RECOMMENDED ⭐

**No API key needed!** Back4app uses GitHub integration (guided setup during deployment).

1. **Sign up**: [https://www.back4app.com/signup](https://www.back4app.com/signup)
2. **Connect GitHub**: Done during `gforge deploy --provider=back4app`

**Free Tier**: 25k container hours/month, 1GB RAM, shared CPU

---

## 🚀 Step 4: Deploy to Production

Gothic Forge supports multiple deployment paths. We recommend starting with **Option A** (full stack).

### **Option A: Full Stack Deploy** (Recommended)

Deploy everything with one command:

```bash
gforge deploy --provider=back4app --with-valkey --with-pages
```

**What This Does:**

1. ✅ **Provisions CockroachDB** - Serverless database cluster
2. ✅ **Provisions Aiven Valkey** - Redis cache instance
3. ✅ **Deploys to Back4app** - Containerized Go app
4. ✅ **Deploys to Cloudflare Pages** - Static assets on CDN

The CLI will guide you through each step with clear prompts.

#### Interactive Prompts

You'll be asked for:

- `SITE_BASE_URL` - Your production URL (e.g., `https://myapp.com`)
- Database provider (CockroachDB recommended, Neon fallback)
- Cache provider (Valkey recommended)
- Cloudflare project name

**Example Session:**

```bash
$ gforge deploy --provider=back4app --with-valkey --with-pages

Gothic Forge v3 :: CLI
Deploy wizard - Provider: back4app

• Checking secrets:
  - Database provider: MISSING (need COCKROACH_API_KEY or NEON_TOKEN)
  - AIVEN_TOKEN: MISSING
  - CLOUDFLARE_API_TOKEN: MISSING

• Enter SITE_BASE_URL: https://myapp.b4a.app
• CockroachDB API key: https://cockroachlabs.cloud/account/api-access
• Enter COCKROACH_API_KEY: <paste-your-key>

✅ CockroachDB: Provisioning serverless cluster...
✅ CockroachDB: Cluster ready (gothic-forge-8x2k)
✅ DATABASE_URL saved to .env

• Aiven token: https://console.aiven.io/account/tokens
• Enter AIVEN_TOKEN: <paste-your-token>

✅ Aiven Valkey: Provisioning Redis instance...
✅ Aiven Valkey: Instance ready (valkey-gothic-forge)
✅ REDIS_URL saved to .env

• Cloudflare API token: https://dash.cloudflare.com/profile/api-tokens
• Enter CLOUDFLARE_API_TOKEN: <paste-your-token>

✅ Building static export...
✅ Cloudflare Pages: Deploying to gothic-forge-demo...
✅ Deployed: https://gothic-forge-demo.pages.dev

• Back4app: Guided container setup

Step 1/5: Create Back4app Account
→ Visit: https://www.back4app.com/signup
→ Sign up and verify your email
Ready to continue? [y/N]: y

Step 2/5: Connect GitHub Repository
→ Visit: https://dashboard.back4app.com/apps
→ Create new "Container App"
→ Connect your GitHub repository
→ Branch: main
Ready to continue? [y/N]: y

Step 3/5: Configure Environment Variables
→ In Back4app dashboard, add these env vars:
  DATABASE_URL=<shown-in-terminal>
  REDIS_URL=<shown-in-terminal>
  JWT_SECRET=<shown-in-terminal>
  SITE_BASE_URL=https://myapp.b4a.app
Ready to continue? [y/N]: y

Step 4/5: Dockerfile Detection
✅ Dockerfile found
✅ .dockerignore found
→ Back4app will auto-build from your Dockerfile

Step 5/5: Deploy!
→ In Back4app, deploy your app
→ Wait for build to complete (~2-3 minutes)
→ Copy your deployment URL

• Enter your Back4app app URL: https://myapp.b4a.app
✅ B4A_APP_URL saved to .env

────────────────────────────────────────
✅ Deployment Complete!

Your app is live at:
  • Go App: https://myapp.b4a.app
  • Static Assets: https://gothic-forge-demo.pages.dev
  • Database: CockroachDB Serverless (gothic-forge-8x2k)
  • Cache: Aiven Valkey (valkey-gothic-forge)

Health checks:
  • https://myapp.b4a.app/healthz
  • https://myapp.b4a.app/readyz

Next steps:
  1. Set up custom domain (optional)
  2. Configure GitHub Actions for CI/CD
  3. Monitor logs: gforge logs
```

---

### **Option B: Static Site Only** (Fastest)

Deploy just the static site (no backend):

```bash
# Generate static HTML
gforge export

# Deploy to Cloudflare Pages
gforge deploy pages --project=gothic-forge-demo

# Or combine:
gforge deploy pages --project=gothic-forge-demo --run
```

**Best For**: Marketing sites, documentation, landing pages

---

### **Option C: Dry-Run First** (Safest)

Preview what will happen without making any changes:

```bash
# See what the deploy wizard would do
gforge deploy --dry-run

# Check prerequisites
gforge deploy --check
```

**Best For**: First-time users, testing configuration

---

## 🔍 Step 5: Verify Deployment

### 5.1 Check Health Endpoints

```bash
# Test your deployed app
curl https://myapp.b4a.app/healthz
# Expected: OK

curl https://myapp.b4a.app/readyz
# Expected: db: OK
#           valkey: OK
#           ready
```

### 5.2 Test Database Connection

```bash
# Run migrations (if not auto-run)
gforge db --migrate

# Check migration status
gforge db --status
```

### 5.3 View Logs

```bash
# Tail application logs
gforge logs

# Or provider-specific
# Back4app: Check dashboard
# Railway: railway logs
```

---

## 🛠️ Troubleshooting

### Problem: "Docker not found"

**For Back4app deployments only**

**Solution**:
1. **Install Docker Desktop**: [Windows](https://docs.docker.com/desktop/install/windows-install/) | [macOS](https://docs.docker.com/desktop/install/mac-install/) | [Linux](https://docs.docker.com/desktop/install/linux-install/)
2. **Restart terminal** after install
3. **Verify**: `docker --version`

**Note**: Docker is only needed for local testing. Back4app builds in the cloud.

### Problem: "COCKROACH_API_KEY: invalid"

**Solution**:
1. Verify key from [API Access page](https://cockroachlabs.cloud/account/api-access)
2. Ensure no extra spaces: `gforge secrets --set COCKROACH_API_KEY=<key>`
3. Key should be ~200+ characters long

### Problem: "Cloudflare authentication failed"

**Solution**:
1. Use `CLOUDFLARE_API_TOKEN` not `CF_API_TOKEN`
2. Token needs **"Edit Cloudflare Workers"** permissions
3. Recreate token: [API Tokens](https://dash.cloudflare.com/profile/api-tokens)

### Problem: "Database migrations failed"

**Solution**:
```bash
# Check connection
gforge db --ping

# Reset and retry
gforge db --reset
gforge db --migrate
```

### Problem: Missing Tools

**Solution**:
```bash
# Auto-fix common issues
gforge doctor --fix

# Manual install if needed
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/gotailwindcss/tailwind/cmd/tailwindcss@latest
```

---

## 📚 Next Steps

### Set Up Custom Domain

**Cloudflare Pages:**
1. Go to Workers & Pages → Your Project → Custom domains
2. Add your domain
3. Update DNS records

**Back4app:**
1. Dashboard → Settings → Custom Domain
2. Add domain and verify

### Configure CI/CD

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Install gforge
        run: go install ./cmd/gforge
      
      - name: Build
        run: gforge build
      
      - name: Deploy Pages
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
        run: gforge deploy pages --project=gothic-forge-demo --run
```

### Monitor Performance

```bash
# View logs
gforge logs

# Check health
curl https://myapp.b4a.app/readyz

# Database metrics
# CockroachDB: https://cockroachlabs.cloud/clusters
# Aiven: https://console.aiven.io
```

### Add More Features

```bash
# Scaffold CRUD resource
gforge add resource Product

# Add OAuth (GitHub)
gforge add oauth github

# Add database table
gforge add db table users
```

---

## 🎓 Learning Resources

- **Gothic Forge Docs**: [Your docs URL]
- **Philosophy**: Read `PHILOSOPHY.md` for design decisions
- **Test Findings**: See `TEST_FINDINGS.md` for known issues
- **Examples**: Check `app/` for templates and patterns

---

## 🆘 Getting Help

- **GitHub Issues**: [Your repo issues]
- **Discussions**: [Your repo discussions]
- **Discord**: [Your Discord invite]

---

## ⚡ Quick Reference

```bash
# Development
gforge dev              # Start dev server
gforge build            # Build production binary
gforge test             # Run tests

# Secrets
gforge secrets --gen-jwt                    # Generate JWT
gforge secrets --set KEY=value              # Set secret
gforge secrets --get KEY                    # Get secret

# Database
gforge db --migrate     # Run migrations
gforge db --status      # Check status
gforge db --reset       # Reset database

# Deployment
gforge deploy --check                       # Preflight check
gforge deploy --dry-run                     # Preview
gforge deploy --provider=back4app           # Deploy to Back4app
gforge deploy pages --project=NAME          # Deploy static site

# Maintenance
gforge doctor           # Health check
gforge doctor --fix     # Auto-fix issues
gforge logs             # View logs
```

---

**🎉 Congratulations!** You've deployed your first Gothic Forge app to production!

Your stack:
- ✅ Go app on Back4app Containers
- ✅ PostgreSQL on CockroachDB Serverless  
- ✅ Redis on Aiven Valkey
- ✅ Static assets on Cloudflare Pages

**All with generous free tiers** - perfect for side projects and MVPs.

Now build something amazing! 🚀
