# Gothic Forge v3

Lean, batteries-included Go starter. Simple, fast, secure, and great DX.

Gothic Forge v3 is built with Go, [chi](https://github.com/go-chi/chi),
[Templ](https://github.com/a-h/templ), [HTMX](https://htmx.org/),
[`gotailwindcss`](https://github.com/gotailwindcss/tailwind) (pure‑Go Tailwind), and
[DaisyUI](https://daisyui.com/). It ships secure defaults (CSP, CSRF, rate limiting),
server‑side rendering, and a fast developer experience with hot reload.

## Philosophy: Teaching Through Doing

Gothic Forge embraces an **opinionated, educational approach** to web development:

- **Batteries-included** - Sane defaults that work out of the box
- **Guided learning** - Some deployment steps teach you WHY, not just HOW
- **Production-ready** - Omakase stack choices based on real-world experience
- **Developer empowerment** - Learn the platform, don't just use a black box

### Why Some Steps Are Manual

For certain platforms (like Back4app Containers), we intentionally guide you through manual setup steps instead of automating everything. This helps you:

1. **Understand your deployment** - Know exactly what's running where
2. **Debug effectively** - When things go wrong, you know the architecture
3. **Make informed choices** - Learn why we chose these specific services
4. **Gain transferable skills** - These patterns apply beyond Gothic Forge

After the initial guided setup, everything runs automatically via `git push` or CLI commands.

## Stack

- Go
- chi (router + middlewares)
- Templ (type-safe UI)
- HTMX (progressive interactivity)
- gotailwindcss (pure‑Go Tailwind)
- DaisyUI (via CDN)

## Features

- **Secure-by-default middleware**: Request ID, Real IP, Recoverer, CORS, rate limit (`go-chi/httprate`),
  session cookies (`scs`), CSP, and optional CSRF in production.
- **SSR with Templ**: Components in `app/templates/` rendered on the server.
- **Pure Go Tailwind CSS**: No Node required. `gotailwindcss` produces `app/styles/output.css` from
  `app/styles/tailwind.input.css` (or your inputs).
- **Hot reload**: `gforge dev` runs Templ generation/watch, Tailwind build/rebuild, and reloads the server.
- **SEO basics**: Favicon, meta tags (Open Graph, Twitter), `/robots.txt` and `/sitemap.xml`; JSON‑LD via `LayoutSEO`.
  - `SEO_KEYWORDS` env lets you override the default keywords included by `LayoutSEO`.
  - `sitemap.xml` includes `<lastmod>` for all URLs.
- **Clean routing**: `app/routes/routes.go` mounts core routes; per‑page registrars via `RegisterRoute`.
- **Tests UX**: `gforge test` builds the server first and runs the suite, with quiet logs.

## Quick start

Prerequisites:

- **Go 1.22+** (required)
- **Git** (required for deployments)
- **Docker** (required for Back4app Containers deployments)
- Optional CLIs: `templ`, `gotailwindcss` (auto-checked by `gforge doctor`)

Run `gforge doctor --fix` to check all prerequisites and get installation guidance.

Doctor:

```powershell
go run ./cmd/gforge doctor
```

Dev:

```powershell
go run ./cmd/gforge dev
# Open http://127.0.0.1:8080/
```

Test:

```powershell
go run ./cmd/gforge test --with-build
```

Build:

```powershell
go run ./cmd/gforge build
```

## Routes

- `/` — Home (Templ: `templates.Index()`)
- `POST /counter/sync` — Server-side counter sample
- `/favicon.ico` — 301 → `/static/favicon.svg`
- `/robots.txt` — Defaults or stream `app/static/robots.txt`
- `/sitemap.xml` — Defaults or stream `app/static/sitemap.xml`
- `/db/posts` — Sample DB‑backed feature (requires `DATABASE_URL`; POST/PUT/DELETE require JWT)
- `/static/*` — Files under `app/static`
- `/static/styles/*` — Files under `app/styles`

### Health Check Endpoints

Production-grade health monitoring endpoints following Kubernetes best practices:

- **`/healthz`** — Basic health check (always returns `ok` if app is running)
  - Used by: Docker HEALTHCHECK, uptime monitors
  - Returns: 200 OK with `ok` response

- **`/livez`** — Liveness probe (process health check)
  - Used by: Kubernetes liveness probes
  - Returns: 200 OK with `alive` response
  - Purpose: Container should be restarted if this fails

- **`/readyz`** — Readiness probe (dependency health checks)
  - Used by: Kubernetes readiness probes, load balancers
  - Checks: Database connectivity (if configured), Valkey/Redis (if configured)
  - Returns: 200 OK with detailed status when ready
  - Returns: 503 Service Unavailable when dependencies are down
  - Purpose: Remove pod from load balancer rotation if not ready

Example readiness response:
```
valkey: OK
db: OK
ready
```

Main entry: `app/routes/routes.go`.

## Scaffolding

```powershell
go run ./cmd/gforge add page about
# -> app/templates/page_about.go
# -> app/routes/page_about.go

go run ./cmd/gforge add component Card
# -> app/templates/component_card.go

go run ./cmd/gforge add auth
# -> /login, /logout + template

go run ./cmd/gforge add oauth github
# -> /oauth/github/{start,callback}

go run ./cmd/gforge add db appdata
# -> app/db/appdata.sql

go run ./cmd/gforge add module blog
# -> page + db scaffold
```

## Project layout

```
app/
  routes/      # chi routes and registrars
  static/      # static assets (favicon, tailwind inputs, etc.)
  styles/      # generated CSS and overrides (served at /static/styles)
  templates/   # Templ components (pure Go)
cmd/
  gforge/      # CLI (doctor, dev, build, test, add, etc.)
  server/      # main web server entrypoint
internal/
  env/         # env helpers
  execx/       # exec helpers
  server/      # router constructor, middlewares, CSP, static mounting
```

## Environment

Copy `.env.example` to `.env` and set:

```
APP_ENV=development
HTTP_HOST=127.0.0.1
HTTP_PORT=8080
LOG_FORMAT=
CORS_ORIGINS=
SITE_BASE_URL=http://127.0.0.1:8080
SEO_KEYWORDS=
DATABASE_URL=
```

- `LOG_FORMAT`: `json` for JSON logs, `off|silent|none` to disable request logs.
- `CORS_ORIGINS`: comma-separated origins (use `*` in dev only).
- `SITE_BASE_URL`: absolute base used by SEO helpers and generated sitemap links.

## Database & Migrations

Gothic Forge uses **CockroachDB Serverless** as the opinionated database standard, with PostgreSQL compatibility via `pgx` and SQL migrations via `goose`.

### Why CockroachDB Serverless?

- **PostgreSQL-compatible** - Works with existing PostgreSQL tools and libraries
- **True serverless** - Pay only for what you use, scales to zero
- **Global distribution** - Low latency worldwide with automatic replication
- **Built-in resilience** - Automatic failover and high availability
- **Educational value** - Learn distributed SQL and modern cloud-native architecture

### 1) Automatic Provisioning (Recommended)

The `gforge deploy` command automatically provisions and configures your database:

```bash
# Set your CockroachDB API key in .env
COCKROACH_API_KEY=your_api_key_here

# Deploy will automatically:
# 1. Create a serverless cluster
# 2. Configure database and user
# 3. Generate secure connection string
# 4. Run migrations automatically
gforge deploy
```

**Get your API key**: https://cockroachlabs.cloud/signup

### 2) Manual Setup (Alternative)

If you prefer manual setup or want to use an existing cluster:

1. Create a CockroachDB Serverless cluster at https://cockroachlabs.cloud
2. Copy the connection string and set it in `.env`:

```bash
DATABASE_URL=postgresql://<user>:<password>@<host>:26257/<db>?sslmode=verify-full
```

**Note**: CockroachDB uses `sslmode=verify-full` for enhanced security.

### 3) Using Neon (Fallback Option)

If you prefer Neon Postgres, it's fully supported:

```bash
# Set NEON_TOKEN instead of COCKROACH_API_KEY
NEON_TOKEN=your_neon_token_here

# Or manually set DATABASE_URL
DATABASE_URL=postgres://<user>:<password>@<host>.neon.tech/<db>?sslmode=require
```

### 4) Working with Migrations

Migrations are located in `app/db/migrations/` and use the goose format.

#### Create and run migrations

- Create a migration file:

```powershell
go run ./cmd/gforge add migration create_posts
```

- Edit the generated file in `app/db/migrations/` and add SQL, e.g.,

```
-- +goose Up
CREATE TABLE posts (
  id bigserial PRIMARY KEY,
  title text NOT NULL,
  body text NOT NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS posts;
```

- Apply migrations:

```powershell
go run ./cmd/gforge db --migrate
```

- Check status or reset:

```powershell
go run ./cmd/gforge db --status
go run ./cmd/gforge db --reset
```

### 4) Run dev and verify

```powershell
go run ./cmd/gforge dev
```

- `/readyz` → should show `db: OK` when `DATABASE_URL` is set and reachable.
- `/db/posts` → sample list/form UI backed by Postgres.

Notes:
- Mutations under `/db/posts` require a valid `gf_jwt` cookie (JWT). Use your OAuth flow or wire a dev-only login helper if needed.

## Security

- CSP is set per environment. In development, inline script/style is allowed for DX.
  In production, inline style is allowed; scripts are restricted to `self` + known CDNs
  (unpkg/jsDelivr) to support HTMX/Alpine and JSON‑LD where needed.
- CSRF middleware is enabled automatically when `APP_ENV=production`.
- Sessions use secure cookie defaults (`HttpOnly`, `SameSite=Lax`, `Secure` in production).

## CI & Releases

- See `.github/workflows/ci.yml` for vet/test/govulncheck on Windows/macOS/Linux.
- See `.github/workflows/release.yml` + `.goreleaser.yaml` to build `gforge` and `gothic-forge-server`
  on tag push (`v*.*.*`).

## Contributing

See `CONTRIBUTING.md`. Follow Conventional Commits.

## License

MIT — see `LICENSE`.

## Deployment

### Omakase Stack Choices

Gothic Forge supports multiple deployment providers with different philosophies:

| Provider | Approach | Best For | Requires |
|----------|----------|----------|----------|
| **Railway** | Automated CLI | Fast iteration, existing users | Railway CLI, tokens |
| **Back4app** | Guided manual | Learning, GitHub workflow | Git, Docker, GitHub repo |

Choose via `--provider` flag: `gforge deploy --provider=railway` (default) or `--provider=back4app`.

### First Deploy (quick guide)

1) Prepare `.env`:

```powershell
cp .env.example .env
go run ./cmd/gforge secrets --set SITE_BASE_URL=https://your-domain
go run ./cmd/gforge secrets --set JWT_SECRET=$(openssl rand -hex 32)
```

2) Preflight and fix:

```powershell
go run ./cmd/gforge doctor --fix
```

3) Choose your provider and deploy:

```powershell
# Railway (automated CLI workflow)
go run ./cmd/gforge deploy --provider=railway --run

# Back4app Containers (guided setup, teaches Docker + CI/CD)
go run ./cmd/gforge deploy --provider=back4app
```

### Token Checklist

**Compute Providers** (choose one):
- Railway:
  - `RAILWAY_TOKEN` (project token) or `RAILWAY_API_TOKEN` (account/team)
- Back4app:
  - `B4A_APP_URL` (saved automatically after guided setup)

**Shared Services**:
- Neon: `NEON_TOKEN`
- Aiven Valkey: `AIVEN_TOKEN`
- Cloudflare Pages: `CF_API_TOKEN`, `CF_ACCOUNT_ID`, `CF_PROJECT_NAME`
- Optional OAuth: `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET`, `OAUTH_BASE_URL` (defaults to `SITE_BASE_URL`)

Store these in `.env` locally. The deploy wizard can sync the key runtime ones to your provider.

### Back4app Containers (guided, educational)

Back4app Containers uses a **guided manual setup** to teach Docker containerization and GitHub-based CI/CD.

**Why guided instead of automated?**
- You learn Docker containerization workflow
- Understand GitHub Actions integration
- Practice environment variable management
- Gain transferable DevOps knowledge

**Prerequisites**:
- Git installed and configured
- Docker installed and running
- GitHub repository created
- Dockerfile in project root

**First-time setup**:

```powershell
# Interactive guided setup (one-time only)
go run ./cmd/gforge deploy --provider=back4app
```

The wizard will walk you through:
1. Creating Back4app account
2. Connecting your GitHub repository
3. Configuring environment variables
4. Watching the initial Docker build
5. Saving your deployment URL

**Subsequent deployments** (the easy way):

```bash
git commit -am "your changes"
git push origin main
# Back4app auto-builds and deploys! ✨
```

**What you learn**:
- Docker image building and containerization
- GitHub webhooks and auto-deployment
- Environment-based configuration
- Zero-downtime rolling deployments
- Platform debugging and log analysis

**Troubleshooting**:
- Check `gforge doctor` for Git/Docker status
- Ensure Dockerfile exists in project root
- Verify environment variables in Back4app dashboard
- View deployment logs at https://dashboard.back4app.com/apps

### Railway (automated CLI)

Use the deploy wizard to guide environment setup and deploy. It checks required secrets and can run an interactive Railway flow.

```powershell
# dry run (no external calls): shows missing secrets and steps
go run ./cmd/gforge deploy --dry-run

# preflight check (no writes/no external actions): validates tools, tokens, env
go run ./cmd/gforge deploy --check

# interactive wizard (first time):
go run ./cmd/gforge deploy --run

# Flags:
#   --init-project   create/link Railway project if missing (wizard will prompt)
#   --project-name   defaults to gothic-forge-v3
#   --service-name   defaults to web
#   --install-tools  attempt to install Railway CLI if missing
```

Required env (typically stored in `.env` or Railway variables):

```
SITE_BASE_URL=https://your-domain
JWT_SECRET=<generated>

# Optional tokens/keys for provider automation
RAILWAY_TOKEN=...          # project token
RAILWAY_API_TOKEN=...      # account/team token (for create/link)
NEON_TOKEN=...
AIVEN_TOKEN=...
CF_API_TOKEN=...
CF_PROJECT_NAME=...
```

### Cloudflare Pages (static)

Export and deploy static HTML to Cloudflare Pages. `_headers` is generated with security/caching defaults.

```powershell
# one-shot deploy with wrangler (if installed)
go run ./cmd/gforge deploy pages --run --project <pages-project-name>

# or dry-run to see the command printed
go run ./cmd/gforge deploy pages --project <pages-project-name>
```

Wrangler install:

Use Homebrew or prebuilt binaries (no Node required):

- macOS: `brew install cloudflare/wrangler/wrangler`
- Windows/Linux: download from https://github.com/cloudflare/wrangler/releases

Notes:
- Export output defaults to `dist/`. Use `--out` to change.
- Security headers (CSP, HSTS, etc.) are written to `dist/_headers`.

### Valkey (Redis-compatible)

Valkey is optional and used for sessions and caching when configured.

Env variables:

```
VALKEY_URL=redis://user:pass@host:port/0
# or REDIS_URL=...
VALKEY_TLS_SKIP_VERIFY=1   # only in dev, if needed
```

`/readyz` will report `valkey: OK|SKIP` automatically.
