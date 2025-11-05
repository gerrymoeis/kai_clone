# Cloudflare Pages Deployment Troubleshooting

## Error 8000000: "An unknown error occurred"

This generic error typically indicates one of the following issues:

### 1. **Free Tier Project Limit Reached** ⭐ Most Common

**Problem**: Cloudflare Free plan allows only **1-3 Pages projects** (varies by account).

**Solution**:
```bash
# List existing projects
wrangler pages project list

# Delete unused projects via Cloudflare Dashboard:
# 1. Go to: https://dash.cloudflare.com/
# 2. Navigate to: Workers & Pages > Pages
# 3. Delete old/unused projects
# 4. Try deployment again
```

**Alternative**: Upgrade to Cloudflare Pro ($20/month) for unlimited projects.

---

### 2. **Account Verification Required**

**Problem**: New Cloudflare accounts may need email/phone verification before creating projects.

**Solution**:
1. Go to Cloudflare Dashboard: https://dash.cloudflare.com/
2. Check for verification prompts (email/phone)
3. Complete verification
4. Wait 5-10 minutes
5. Try deployment again

---

### 3. **API Token Permission Issues**

**Problem**: Your API token may not have sufficient permissions.

**Solution**:

#### **Option A: Use OAuth (Recommended for First Deploy)**
```bash
# Let wrangler handle OAuth automatically
wrangler pages deploy dist --project-name=my-site

# Wrangler will open browser and request all needed permissions
```

#### **Option B: Create New API Token with Correct Scopes**
1. Go to: https://dash.cloudflare.com/profile/api-tokens
2. Click "Create Token"
3. Use template: "Edit Cloudflare Workers"
4. **Add these permissions**:
   - Account > Cloudflare Pages > Edit
   - Zone > Zone > Read (if using custom domains)
5. Click "Continue to summary" → "Create Token"
6. Copy token and add to `.env`:
   ```bash
   CLOUDFLARE_API_TOKEN=your-new-token
   ```

---

### 4. **Project Name Conflicts**

**Problem**: Project name already exists (even in deleted projects, names may be reserved for 24-48 hours).

**Solution**:
```bash
# Try a different project name
./gforge deploy pages --project=my-site-v2 --run

# Or use a unique name with timestamp
./gforge deploy pages --project=my-app-$(date +%s) --run
```

---

### 5. **Rate Limiting**

**Problem**: Too many API requests in short time.

**Solution**:
- Wait 5-10 minutes
- Try again
- Avoid rapid retries

---

### 6. **Account Region/Compliance Issues**

**Problem**: Your Cloudflare account region may have restrictions.

**Solution**:
- Check Cloudflare Dashboard for compliance notices
- Contact Cloudflare support if you're in a restricted region
- Try from a different network/VPN (corporate firewalls may interfere)

---

## Recommended Deployment Flow

### **First-Time Deployment (Use OAuth)**

```bash
# Step 1: Export static site
./gforge export

# Step 2: Deploy with wrangler (OAuth authentication)
wrangler pages deploy dist --project-name=my-site

# Wrangler will:
# - Open browser for OAuth
# - Request all needed permissions
# - Create project
# - Deploy files
```

### **Subsequent Deployments (Use gforge)**

```bash
# After OAuth is set up, use gforge
./gforge deploy pages --project=my-site --run
```

---

## Debugging Steps

### 1. **Check Existing Projects**

```bash
# List all Pages projects
wrangler pages project list

# Output shows:
# - Project names
# - Production domains
# - Creation dates
```

### 2. **Check API Token**

```bash
# Verify token works
wrangler whoami

# Should show:
# - Your email
# - Account name
# - Account ID
```

### 3. **Manual Project Creation**

Create project manually first, then deploy:

```bash
# 1. Create project via Dashboard:
# https://dash.cloudflare.com/ > Workers & Pages > Create

# 2. Deploy to existing project:
wrangler pages deploy dist --project-name=existing-project
```

### 4. **Check Wrangler Logs**

```bash
# Logs are written to:
# Windows: %APPDATA%\xdg.config\.wrangler\logs\
# macOS/Linux: ~/.config/.wrangler/logs/

# View latest log:
# Windows:
dir /O-D "%APPDATA%\xdg.config\.wrangler\logs"

# macOS/Linux:
ls -lt ~/.config/.wrangler/logs/ | head -1
```

---

## Common Error Patterns

| Error Code | Meaning | Solution |
|------------|---------|----------|
| **8000000** | Generic API error | Check limits, verification, token |
| **9000** | Authentication failed | Re-run `wrangler login` |
| **10000** | Rate limited | Wait 5-10 min |
| **81000** | Project limit reached | Delete old projects or upgrade |

---

## Cloudflare Free Tier Limits

| Resource | Free Tier Limit |
|----------|----------------|
| **Pages Projects** | 1-3 projects (varies) |
| **Deployments/month** | 500 |
| **Builds/month** | 500 |
| **Bandwidth** | Unlimited |
| **Custom Domains** | 100 per project |

---

## Alternative: Direct Upload (Without Project Creation)

If you have an existing project:

```bash
# Deploy to existing project (no creation needed)
wrangler pages deploy dist --project-name=existing-project

# Wrangler will:
# - Skip project creation
# - Upload files directly
# - Create new deployment
```

---

## Still Getting Errors?

### **Option 1: Use Cloudflare Dashboard**

Manual deployment via dashboard:

1. Go to: https://dash.cloudflare.com/
2. Workers & Pages > Create > Pages
3. Upload folder: Select `dist/` directory
4. Configure: Set project name, branch
5. Deploy

### **Option 2: Contact Cloudflare Support**

If none of the above work:

1. Go to: https://cfl.re/3WgEyrH
2. Include:
   - Error code: 8000000
   - Account ID (from `wrangler whoami`)
   - Timestamp of error
   - Wrangler log file path

---

## Prevention Tips

1. **Use OAuth for first deploy** - Ensures correct permissions
2. **Check project limits before deployment** - Run `wrangler pages project list`
3. **Use descriptive project names** - Easier to manage
4. **Clean up old projects** - Stay within limits
5. **Verify account immediately** - Avoid delays

---

## Quick Fix Checklist

- [ ] Check existing projects: `wrangler pages project list`
- [ ] Verify account: Check Cloudflare dashboard for notices
- [ ] Try different project name
- [ ] Wait 5-10 minutes if rate limited
- [ ] Use OAuth: `wrangler pages deploy dist --project-name=test`
- [ ] Check API token permissions
- [ ] Delete old projects if at limit
- [ ] Contact Cloudflare support if persistent

---

**Most likely issue**: **Project limit reached on free tier**. Delete old projects or upgrade plan.
