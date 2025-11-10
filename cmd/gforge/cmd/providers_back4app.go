package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gothicforge3/internal/execx"
)

// back4appGuidedSetup walks the user through setting up Back4app Containers.
// This is intentionally a guided (not automated) flow to teach the platform and
// encourage learning. Returns error if setup is cancelled or fails validation.
func back4appGuidedSetup(ctx context.Context, dryRun bool) error {
	if dryRun {
		fmt.Println("  • Back4app Containers (dry-run): would guide manual setup")
		return nil
	}

	// Check if already deployed
	if url := strings.TrimSpace(os.Getenv("B4A_APP_URL")); url != "" {
		fmt.Printf("  • Back4app: already deployed at %s\n", url)
		fmt.Println("    → Future deploys: git push (auto-deploy enabled)")
		return nil
	}

	// Verify Git is available (essential for GitHub integration)
	if _, ok := execx.Look("git"); !ok {
		fmt.Println("  • Back4app Containers requires Git for GitHub integration")
		printGitInstallHelp()
		return fmt.Errorf("git not found; install Git and re-run")
	}

	// Verify Docker is available
	if _, ok := execx.Look("docker"); !ok {
		fmt.Println("  • Back4app Containers requires Docker")
		printDockerInstallHelp()
		return fmt.Errorf("docker not found; install Docker and re-run")
	}

	// Check Docker daemon
	if err := execx.Run(ctx, "docker ps", "docker", "ps"); err != nil {
		fmt.Println("  • Docker daemon not running")
		fmt.Println("    → Start Docker Desktop and re-run gforge deploy")
		return fmt.Errorf("docker daemon unavailable")
	}

	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("  GUIDED SETUP: Back4app Containers")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("Back4app Containers offers seamless Docker deployment:")
	fmt.Println("  ✓ Auto-deploy on git push (after initial setup)")
	fmt.Println("  ✓ Free tier with generous limits")
	fmt.Println("  ✓ Zero downtime deployments")
	fmt.Println("  ✓ Built-in monitoring and logs")
	fmt.Println()
	fmt.Println("WHY GUIDED SETUP?")
	fmt.Println("This setup is intentionally guided (not automated) so you:")
	fmt.Println("  • Understand how your app is deployed")
	fmt.Println("  • Can troubleshoot issues independently")
	fmt.Println("  • Learn the Back4app platform deeply")
	fmt.Println("  • Gain transferable DevOps knowledge")
	fmt.Println()
	fmt.Println("Gothic Forge Philosophy: Teaching Through Doing")
	fmt.Println("We believe in empowering developers, not creating black boxes.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ready to proceed? [Y/n]: ")
	ans, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(ans)) == "n" {
		return fmt.Errorf("setup cancelled by user")
	}

	return printBack4appSteps(reader)
}

// printBack4appSteps displays step-by-step instructions with validation checkpoints.
func printBack4appSteps(reader *bufio.Reader) error {
	// Step 1: Account Setup
	fmt.Println("\n──────────────────────────────────────────")
	fmt.Println("STEP 1: Back4app Account")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("1. Open: https://www.back4app.com/signup")
	fmt.Println("2. Sign up (or login if you have an account)")
	fmt.Println("3. Verify your email address")
	fmt.Println()
	fmt.Println("WHY THIS MATTERS:")
	fmt.Println("  Back4app offers a generous free tier for containers,")
	fmt.Println("  making it perfect for learning and side projects.")
	fmt.Println()
	fmt.Print("✓ Account ready? Press ENTER to continue...")
	_, _ = reader.ReadString('\n')

	// Step 2: Connect GitHub - WITH CRITICAL PRODUCT SELECTION WARNING
	fmt.Println("\n──────────────────────────────────────────")
	fmt.Println("STEP 2: Connect GitHub Repository")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println()
	fmt.Println("⚠️  CRITICAL: You MUST select the CORRECT product!")
	fmt.Println()
	fmt.Println("Back4app has TWO different products:")
	fmt.Println()
	fmt.Println("  ❌ Parse Server (Backend-as-a-Service)")
	fmt.Println("     - MongoDB database")
	fmt.Println("     - Parse API")
	fmt.Println("     - NOT what we need!")
	fmt.Println()
	fmt.Println("  ✅ Containers (Docker Hosting)")
	fmt.Println("     - Deploy Docker images")
	fmt.Println("     - Use your own database")
	fmt.Println("     - THIS is what Gothic Forge needs!")
	fmt.Println()
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("DEPLOYMENT STEPS:")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("1. Go to: https://dashboard.back4app.com")
	fmt.Println("2. Click 'Build new app'")
	fmt.Println("3. Look for TWO tabs at the top:")
	fmt.Println("   - 'Backend as a Service' (Parse Server)")  
	fmt.Println("   - 'Container as a Service'")
	fmt.Println()
	fmt.Println("4. Click 'Container as a Service' tab ← IMPORTANT!")
	fmt.Println("5. Click 'GitHub' under deployment options")
	fmt.Println("6. Authorize Back4app to access your repositories")
	fmt.Println("7. Select your repository")
	fmt.Println()

	// Auto-detect git remote
	if repoURL := detectGitRemote(); repoURL != "" {
		fmt.Printf("   Detected repository: %s\n", repoURL)
		fmt.Println()
	}

	fmt.Println("──────────────────────────────────────────")
	fmt.Println("VERIFICATION:")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("After creating the app, verify you see:")
	fmt.Println("  ✅ 'Container' in the app type")
	fmt.Println("  ✅ Option to configure Dockerfile")
	fmt.Println()
	fmt.Println("If you see:")
	fmt.Println("  ❌ 'Parse Server Version'")
	fmt.Println("  ❌ 'Database: MongoDB'")
	fmt.Println("  ❌ 'API URL: parseapi.back4app.com'")
	fmt.Println()
	fmt.Println("  → You created the WRONG type! DELETE it and start over.")
	fmt.Println()
	fmt.Println("WHY THIS MATTERS:")
	fmt.Println("  Parse Server uses MongoDB and a fixed API structure.")
	fmt.Println("  Gothic Forge needs full control with Docker containers")
	fmt.Println("  to connect to CockroachDB and deploy custom Go apps.")
	fmt.Println()
	fmt.Print("✓ Container app created (not Parse Server)? Press ENTER to continue...")
	_, _ = reader.ReadString('\n')

	// Step 3: Configure Environment
	fmt.Println("\n──────────────────────────────────────────")
	fmt.Println("STEP 3: Configure Deployment")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("App Configuration:")
	fmt.Println("  • App Name: gothic-forge-app (or your choice)")
	fmt.Println("  • Branch: main")
	fmt.Println("  • Root Directory: . (leave as default)")
	fmt.Println("  • Dockerfile: Dockerfile (auto-detected)")
	fmt.Println()
	fmt.Println("IMPORTANT: Ensure you have a Dockerfile in your repo root.")
	fmt.Println("If missing, Back4app will fail to build.")
	fmt.Println()

	// Check if Dockerfile exists
	if _, err := os.Stat("Dockerfile"); os.IsNotExist(err) {
		fmt.Println("⚠️  WARNING: Dockerfile not found in current directory!")
		fmt.Println("    Create one before proceeding or deployment will fail.")
		fmt.Println()
		fmt.Print("Continue anyway? [y/N]: ")
		ans, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(ans)) != "y" {
			return fmt.Errorf("setup cancelled: missing Dockerfile")
		}
	} else {
		fmt.Println("✓ Dockerfile found in repository")
	}
	fmt.Println()

	fmt.Println("Environment Variables (click 'Add Variable' in Back4app UI):")
	fmt.Println("Copy these from your .env file:")
	fmt.Println()

	kv := loadEnvFile(".env")
	requiredVars := []string{"DATABASE_URL", "JWT_SECRET", "SITE_BASE_URL"}
	optionalVars := []string{"REDIS_URL", "VALKEY_URL", "GITHUB_CLIENT_ID", "GITHUB_CLIENT_SECRET"}

	fmt.Println("Required:")
	for _, key := range requiredVars {
		val := kv[key]
		if val != "" {
			fmt.Printf("  • %s = %s\n", key, maskValue(val))
		} else {
			fmt.Printf("  • %s = (⚠️  MISSING - set this in Back4app UI)\n", key)
		}
	}

	fmt.Println("\nOptional:")
	for _, key := range optionalVars {
		val := kv[key]
		if val != "" {
			fmt.Printf("  • %s = %s\n", key, maskValue(val))
		}
	}

	fmt.Println()
	fmt.Println("WHY ENVIRONMENT VARIABLES:")
	fmt.Println("  These configure your app at runtime without hardcoding")
	fmt.Println("  secrets in your codebase. This is a security best practice.")
	fmt.Println()
	fmt.Print("✓ Environment configured? Press ENTER to continue...")
	_, _ = reader.ReadString('\n')

	// Step 4: Deploy
	fmt.Println("\n──────────────────────────────────────────")
	fmt.Println("STEP 4: Initial Deployment")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("1. Click 'Create App' in Back4app UI")
	fmt.Println("2. Back4app will:")
	fmt.Println("   • Clone your repository")
	fmt.Println("   • Build your Docker image")
	fmt.Println("   • Deploy to a container")
	fmt.Println("   • Provide a live URL: <app-name>.b4a.run")
	fmt.Println()
	fmt.Println("3. Wait for build to complete (typically 2-5 minutes)")
	fmt.Println("4. Check deployment logs for any errors")
	fmt.Println()
	fmt.Println("WHY THIS TAKES TIME:")
	fmt.Println("  Docker builds your entire app from scratch, ensuring")
	fmt.Println("  consistency between dev and production environments.")
	fmt.Println()
	fmt.Print("✓ Deployment complete? Press ENTER to continue...")
	_, _ = reader.ReadString('\n')

	// Step 5: Save deployment URL
	fmt.Println("\n──────────────────────────────────────────")
	fmt.Println("STEP 5: Save Deployment URL")
	fmt.Println("──────────────────────────────────────────")
	fmt.Print("Enter your Back4app app URL (e.g., myapp.b4a.run): ")
	appURL, _ := reader.ReadString('\n')
	appURL = strings.TrimSpace(appURL)

	if appURL != "" {
		// Normalize URL
		appURL = strings.TrimPrefix(appURL, "http://")
		appURL = strings.TrimPrefix(appURL, "https://")
		fullURL := "https://" + appURL

		// Save to .env
		kv := map[string]string{"B4A_APP_URL": fullURL}
		if err := updateEnvFileInPlace(".env", kv); err == nil {
			fmt.Println("  → B4A_APP_URL saved to .env")
		} else {
			fmt.Println("  → Warning: failed to save to .env:", err)
		}
	}

	// Success summary
	fmt.Println("\n════════════════════════════════════════════════════════════")
	fmt.Println("  ✓ Back4app Containers Setup Complete!")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("\n🎉 Congratulations! Your app is now live at:")
	if appURL != "" {
		fmt.Printf("   https://%s\n", appURL)
	}
	fmt.Println()
	fmt.Println("Next Deployment (The Easy Way):")
	fmt.Println("  1. Make code changes locally")
	fmt.Println("  2. git commit -am \"your changes\"")
	fmt.Println("  3. git push origin main")
	fmt.Println("  4. Back4app auto-deploys! ✨")
	fmt.Println()
	fmt.Println("Monitor your app:")
	fmt.Println("  • Dashboard: https://dashboard.back4app.com/apps")
	fmt.Println("  • View logs, metrics, and manage env variables")
	fmt.Println()
	fmt.Println("What You Learned:")
	fmt.Println("  ✓ Docker containerization workflow")
	fmt.Println("  ✓ GitHub-based CI/CD pipeline")
	fmt.Println("  ✓ Environment variable management")
	fmt.Println("  ✓ Zero-downtime deployment strategies")
	fmt.Println()

	return nil
}

// detectGitRemote attempts to detect the current repository's remote URL.
// Returns empty string if git is not available or no remote is configured.
func detectGitRemote() string {
	// Check if git is available
	if _, ok := execx.Look("git"); !ok {
		return ""
	}
	
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	url := strings.TrimSpace(string(out))
	// Clean up SSH URLs to HTTPS format for display
	url = strings.Replace(url, "git@github.com:", "https://github.com/", 1)
	url = strings.TrimSuffix(url, ".git")
	return url
}

// maskValue masks sensitive values for display, showing only first/last 4 chars.
func maskValue(val string) string {
	if val == "" {
		return "(empty)"
	}
	if len(val) <= 8 {
		return "***"
	}
	return val[:4] + "***" + val[len(val)-4:]
}
