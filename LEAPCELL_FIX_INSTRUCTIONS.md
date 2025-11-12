# 🔧 Fix Leapcell Deployment - CSS 404 Errors

**Problem**: CSS files return 404 Not Found  
**Root Cause**: CSS files weren't being built during deployment  
**Solution**: Use `build.sh` script that handles all build steps

---

## ✅ What You Need To Do

### **Step 1: Update Build Command in Leapcell Dashboard**

1. **Go to Leapcell Dashboard**:
   - Visit: https://leapcell.io/dashboard
   - Select: `demo-htmx-v71` service

2. **Go to Settings**:
   - Click on "Settings" tab
   - Find "Build & Deploy" or "Build Settings" section

3. **Update Build Command**:
   - **Old** (WRONG): `go build -o server ./cmd/server`
   - **New** (CORRECT): `bash build.sh`

4. **Verify Start Command**:
   - Should be: `./server`
   - Port: `8080`

5. **Save Settings**

---

### **Step 2: Trigger Manual Redeploy**

Since you updated the build command, you need to manually redeploy:

1. **In Leapcell Dashboard**:
   - Go to "Deployments" tab
   - Click "Redeploy" or "Deploy Now"
   - Or: Click "Trigger Deployment" button

2. **Watch Build Logs**:
   - You should see:
     ```
     🔧 Gothic Forge Build Script for Leapcell
     ==========================================
     📦 Installing build tools...
     🎨 Generating templ templates...
     💅 Building CSS with Tailwind...
     🚀 Building Go server...
     ✅ Build complete! Binary: ./server
     ```

3. **Wait for Deployment**:
   - Build time: ~2-3 minutes
   - Status should change to: "Running" or "Active"

---

### **Step 3: Verify CSS Files Load**

Once deployment completes:

1. **Visit Your Site**:
   ```
   https://demo-htmx-v71-gerrymoeis1981-i2otkmzb.apn.leapcell.dev
   ```

2. **Open Browser Console** (F12):
   - Should see **NO CSS errors**
   - ✅ No more 404 errors for `output.css` or `overrides.css`

3. **Check Styling**:
   - ✅ Page should have proper colors and layout
   - ✅ Tailwind CSS should be working
   - ✅ DaisyUI components styled correctly

---

## 📊 What Changed?

### **Previous Setup** ❌
```bash
Build Command: go build -o server ./cmd/server
```

**Problems**:
- Only built Go binary
- Didn't generate CSS files
- CSS files are gitignored (not in repo)
- Result: 404 errors for CSS

### **New Setup** ✅
```bash
Build Command: bash build.sh
```

**What build.sh does**:
1. ✅ Installs `templ` tool (for templates)
2. ✅ Installs `gotailwindcss` tool (for CSS)
3. ✅ Generates templ templates
4. ✅ **Builds CSS files** (output.css + minified)
5. ✅ Compiles Go binary

**Result**: CSS files are generated during build → No more 404!

---

## 🔍 Troubleshooting

### **If Build Fails**

**Check Build Logs** for these common issues:

1. **"bash: command not found"**
   - Leapcell should have bash installed
   - Try: Runtime = "Go" (not "Docker")

2. **"templ: command not found"**
   - The script installs it, but PATH might be wrong
   - Check if `$(go env GOPATH)/bin` is added to PATH

3. **"gotailwindcss: command not found"**
   - Same as above
   - Script adds tools to PATH

4. **"./server: not found" (during start)**
   - Build failed silently
   - Check if `go build` step completed
   - Look for errors in build logs

### **If CSS Still 404s After Successful Build**

1. **Check if files were created**:
   - Build logs should show: "✅ Build complete!"
   - Check deployment artifacts for `app/styles/output.css`

2. **Check server logs**:
   - Does server start successfully?
   - Look for errors loading static files

3. **Check file paths**:
   - Server expects: `/static/styles/output.css`
   - File should be at: `app/styles/output.css`
   - Server mounts `app/styles` at `/static/styles`

---

## 📝 Build Settings Summary

Copy these settings into Leapcell Dashboard:

```yaml
Runtime: Go
Build Command: bash build.sh
Start Command: ./server
Port: 8080
```

---

## 🎯 Expected Result

After following these steps:

✅ **Deployment succeeds**  
✅ **CSS files load correctly** (no 404)  
✅ **Styles display properly**  
✅ **HTMX interactions work**  
✅ **Alpine.js animations work**  

---

## 💡 Why This Happened

**Gothic Forge uses a build pipeline**:
```
Source Files → templ generate → Go templates
              → gotailwindcss → CSS files
              → go build → Binary
```

**The issue**: Original Leapcell setup only did the last step (`go build`), skipping template and CSS generation.

**The fix**: `build.sh` runs the FULL pipeline, ensuring all assets are built before the Go binary.

---

## 📚 Additional Notes

### **About build.sh**

- **Location**: Root of repository
- **Permissions**: Executable by default on Unix systems
- **Exit on error**: Uses `set -e` to fail fast
- **Logging**: Shows progress for each step
- **Path handling**: Adds Go bin to PATH for tools

### **Files Generated During Build**

1. `app/templates/*_templ.go` - Generated Go code from templ templates
2. `app/styles/output.css` - Compiled and minified Tailwind CSS
3. `server` - Compiled Go binary

### **Why CSS is Gitignored**

Gothic Forge follows best practices:
- Generated files should NOT be committed
- They're built during deployment
- Keeps repository clean
- Ensures consistent builds

---

## ✅ Checklist

Before testing:

- [ ] Updated Build Command to `bash build.sh`
- [ ] Verified Start Command is `./server`
- [ ] Saved settings in Leapcell
- [ ] Triggered manual redeploy
- [ ] Watched build logs for success
- [ ] Waited for "Running" status

After deployment:

- [ ] Visited site URL
- [ ] Opened browser console (F12)
- [ ] Verified NO CSS 404 errors
- [ ] Checked page styling works
- [ ] Tested HTMX interactions

---

## 🚀 You're All Set!

Once you update the build command and redeploy, your Gothic Forge app will be fully functional on Leapcell with proper CSS styling! 🎉

**Need Help?**
- Leapcell Docs: https://docs.leapcell.io/
- Leapcell Discord: https://discord.gg/qF7efny8x2
- Gothic Forge Issues: https://github.com/yourusername/gothic-forge/issues
