# Gothic Forge v7.1 - Release Notes

**Release Date**: November 11, 2025  
**Branch**: `stable_v7.1`  
**Status**: ✅ Ready for Production

---

## 🎯 What's New in v7.1

### **Leapcell Integration** (Primary Compute Provider)

Gothic Forge v7.1 introduces **Leapcell** as the default and recommended compute provider, offering superior free tier benefits and alignment with our "teaching through doing" philosophy.

#### **Why Leapcell?**

- 🎁 **20 FREE projects** (vs 1-2 on alternatives)
- 🗄️ **1 FREE PostgreSQL database** per project
- ⚡ **Serverless-first** architecture (pay-per-use)
- 🔄 **Auto-deploy** on git push (GitHub integration)
- 🌍 **Global CDN** included
- 📊 **Built-in monitoring** and analytics
- 🛡️ **DDoS protection** and WAF

---

## 📦 What's Included

### **New Provider: Leapcell**

```bash
# Deploy to Leapcell (default)
gforge deploy

# Explicit provider selection
gforge deploy --provider=leapcell
gforge deploy --provider=back4app    # Alternative
gforge deploy --provider=railway     # Alternative

# Dry-run preview
gforge deploy --dry-run --provider=leapcell
```

### **Features**

1. **Guided Deployment Wizard**
   - Step-by-step instructions
   - GitHub integration guidance
   - Build configuration assistance
   - Deployment URL tracking
   - Environment variable checklist

2. **CSS MIME Type Fix** (Critical)
   - Resolves platform-specific MIME type issues
   - Explicit Content-Type headers for all static assets
   - Cache headers for CDN optimization
   - Ensures styles load correctly on Leapcell

3. **Enhanced Post-Deployment Guidance**
   - Clear environment variable setup instructions
   - Database configuration guide (Leapcell PostgreSQL)
   - Cache setup guide (Aiven Valkey)
   - Troubleshooting tips
   - Monitoring and debugging guidance

---

## 🏗️ Architecture Decision: Monolith + HTMX

### **Why NOT API-Only Backend?**

After deployment testing and user feedback, we confirmed that Gothic Forge's monolith architecture is **correct and intentional**:

**HTMX Requires Server-Rendered HTML**:
```
User Action → HTMX Request → Server renders HTML fragment → HTMX swaps DOM
```

**API-Only Would Break This**:
- ❌ Backend returns JSON instead of HTML
- ❌ Frontend needs React/Vue to render HTML from JSON
- ❌ Defeats entire purpose of Gothic Forge
- ❌ Adds complexity we're avoiding

**Correct Architecture** (v7.1):
```
┌────────────────────────────────────────┐
│  Cloudflare (Optional CDN Proxy)       │
│  - Caches static assets                │
│  - DDoS protection                     │
│  - Global edge network                 │
└──────────────┬─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────┐
│  Leapcell/Railway/Back4app (Origin)    │
│  ────────────────────────────────      │
│  [Go Monolith - Full Stack]            │
│  • Renders HTML templates (templ)      │
│  • Serves static files (/static/*)     │
│  • HTMX endpoints (HTML fragments)     │
│  • Database connection                 │
│  • Cache connection (Valkey)           │
└────────────────────────────────────────┘
```

**Benefits**:
- ✅ HTMX works perfectly
- ✅ No JavaScript framework needed
- ✅ Simple deployment (one codebase)
- ✅ Optional CDN for performance

---

## 🐛 Bugs Fixed

### **1. CSS MIME Type Error** (CRITICAL)

**Issue**: Leapcell served CSS files with `Content-Type: text/plain`, causing browsers to reject stylesheets.

**Symptom**:
```
Refused to apply style from '...output.css' because its MIME type 
('text/plain') is not a supported stylesheet MIME type
```

**Fix**: 
- Updated `internal/server/server.go` - `mountStatic()` function
- Added explicit Content-Type headers for all static file types
- Added Cache-Control headers for CDN optimization

**Impact**: Sites now display correctly with proper styling

---

### **2. Missing Post-Deployment Guidance**

**Issue**: Users unsure what to do after initial deployment - environment variables not configured.

**Fix**:
- Enhanced success message with clear next steps
- Added environment variables checklist
- Database setup guidance (Leapcell PostgreSQL)
- Troubleshooting section

**Impact**: Smoother deployment experience, fewer configuration errors

---

## 📋 Provider Hierarchy (v7.1)

1. **Leapcell** (PRIMARY - Default)
   - Best free tier (20 projects)
   - Includes free PostgreSQL
   - GitHub auto-deploy
   - Recommended for new users

2. **Back4app** (Alternative)
   - Docker-focused
   - GitHub integration
   - Good for Docker-experienced users

3. **Railway** (Alternative)
   - CLI-automated
   - Fastest deployment
   - Good for quick prototypes

---

## 🚀 Deployment Guide

### **Quick Start**

```bash
# 1. Deploy to Leapcell
gforge deploy

# 2. Follow the interactive wizard
#    - Confirm account creation
#    - Connect GitHub repository
#    - Select repository and branch
#    - Configure build settings (provided by wizard)
#    - Save deployment URL

# 3. Configure environment variables in Leapcell dashboard
#    Go to: https://leapcell.io/dashboard
#    → Service Settings → Environment Variables
#    
#    Required:
#    APP_ENV=production
#    SITE_BASE_URL=https://your-app.leapcell.dev
#    JWT_SECRET=<from-local-env>
#    
#    Optional:
#    DATABASE_URL=<from-leapcell-database>
#    VALKEY_URL=<from-aiven-console>

# 4. Click "Redeploy" in Leapcell dashboard

# 5. Future deploys are automatic:
git push origin main  # Auto-deploys to Leapcell!
```

### **Environment Variables Setup**

See detailed guide in deployment success message or:
- **Leapcell Dashboard**: https://leapcell.io/dashboard
- **Leapcell Database**: Free PostgreSQL in Database section
- **Aiven Valkey**: https://console.aiven.io/signup (free tier)

---

## 🔒 Security

### **Built-In Protections**

1. **Rate Limiting** (120 req/min per IP)
2. **CORS** configuration (restrict origins)
3. **CSRF** protection (production mode)
4. **Content Security Policy** (CSP)

### **Backend URL Exposure**

**Concern**: Leapcell URL is public - can be attacked directly.

**Mitigations**:
- ✅ Rate limiting at application layer
- ✅ CORS restrictions
- ✅ CSRF tokens
- 💡 Optional: Add Cloudflare proxy (hides origin IP)
- 💡 Optional: Leapcell firewall (allow only Cloudflare IPs)

---

## 📊 Testing Results

| Component | Status | Notes |
|-----------|--------|-------|
| Leapcell Wizard | ✅ Pass | Excellent UX, clear instructions |
| GitHub Integration | ✅ Pass | Auto-detected repository |
| Build Process | ✅ Pass | Go app built successfully |
| Deployment | ✅ Pass | App accessible at URL |
| CSS Loading | ✅ Pass | Fixed MIME type issue |
| HTMX Functionality | ✅ Pass | Server-rendered HTML works |
| Auto-Redeploy | ✅ Pass | Git push triggers deployment |

**Test Environment**:
- Test App: Gothic Forge demo
- Deployed URL: `https://gforgev7-vanwickbuddy9375-qxmifsdo.apn.leapcell.dev`
- Platform: Leapcell Hobby tier (free)

---

## 🎓 Lessons Learned

1. **HTMX Architecture Validated**
   - Monolith is correct for HTMX-based apps
   - API-only backend would break HTMX model
   - User correctly identified this requirement

2. **Platform-Specific Issues**
   - Different platforms handle MIME types differently
   - Always test on actual deployment platform
   - Local testing doesn't catch all issues

3. **Deployment UX Gaps**
   - Users need clear post-deployment guidance
   - Environment variable setup should be explicit
   - Database/cache integration needs better workflow

---

## 📝 Documentation Updates

### **New Files**

1. `LEAPCELL_DEPLOYMENT_TEST.md` - Comprehensive test findings
2. `DEPLOYMENT_ARCHITECTURE_v7.md` - Architecture analysis
3. `RELEASE_NOTES_v7.1.md` - This file

### **Updated Files**

1. `.env.example` - Added Leapcell variables, updated to v7.0 stack
2. `cmd/gforge/cmd/deploy.go` - Leapcell provider integration
3. `cmd/gforge/cmd/providers_leapcell.go` - New provider implementation
4. `internal/server/server.go` - CSS MIME type fix
5. `.gitignore` - Added test/research documentation

---

## 🔄 Migration from v6.x

### **Breaking Changes**: ❌ NONE

v7.1 is **100% backward compatible** with v6.x deployments.

### **What Changes for Existing Users?**

- Default provider changes from Railway → Leapcell
- Can still use Railway/Back4app with `--provider` flag
- All existing deployments continue to work

### **Recommended Actions**

1. Update to v7.1: `git pull origin stable_v7.1`
2. Try Leapcell for new projects: `gforge deploy`
3. Keep existing deployments on current providers

---

## 🚧 Known Limitations

1. **Cloudflare Pages Integration**
   - Not recommended for HTMX apps (breaks server-rendering)
   - Use for marketing sites/docs only
   - Main app should use monolith deployment

2. **Database/Cache Setup**
   - Currently manual (dashboard configuration)
   - Future: Auto-provision from CLI

3. **Environment Variables**
   - Must be configured in platform dashboard
   - Future: Sync from local .env file

---

## 🛣️ Roadmap (Future Versions)

### **v7.2** (Planned)
- Auto-provision Leapcell PostgreSQL from CLI
- Environment variable sync command
- Database migration automation

### **v7.3** (Planned)
- Cloudflare DNS/Proxy setup wizard
- Custom domain configuration guide
- Enhanced monitoring dashboard

### **v8.0** (Future)
- Multi-region deployment support
- Built-in observability (OpenTelemetry)
- A/B testing framework

---

## 🙏 Credits

**Testing & Feedback**: User feedback during live deployment testing identified critical issues and validated architectural decisions.

**Key Contributions**:
- Identified CSS MIME type issue
- Validated HTMX monolith architecture
- Highlighted security concerns
- Identified deployment UX gaps

---

## 📚 Resources

- **Leapcell Docs**: https://docs.leapcell.io/
- **Leapcell Signup**: https://leapcell.io/signup (20 free projects!)
- **Leapcell Discord**: https://discord.gg/qF7efny8x2
- **Gothic Forge README**: See project root
- **Test Findings**: See `LEAPCELL_DEPLOYMENT_TEST.md`
- **Architecture**: See `DEPLOYMENT_ARCHITECTURE_v7.md`

---

## ✅ Release Checklist

- [x] Leapcell provider implemented
- [x] CSS MIME type fix applied
- [x] Post-deployment guidance enhanced
- [x] Code compiled successfully
- [x] Architecture documented
- [x] Test findings documented
- [x] Release notes created
- [x] Branch `stable_v7.1` created
- [ ] Pushed to origin (ready when you are!)

---

## 🎉 Summary

Gothic Forge v7.1 brings **Leapcell** as the new default compute provider, offering exceptional free tier benefits (20 projects!) and a streamlined deployment workflow. Critical bug fixes ensure proper CSS loading, and enhanced guidance helps users configure their deployments correctly.

The release also **validates** our architectural decision to keep the monolith + HTMX approach, confirming that API-only backends would break HTMX's server-rendering requirement.

**Ready to deploy?**

```bash
gforge deploy
```

Welcome to the future of Go web development! 🚀
