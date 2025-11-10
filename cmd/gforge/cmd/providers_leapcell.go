package cmd

import (
	"bufio"
	"context"
	"fmt"
	"strings"
)

// runLeapcellDeploy guides the user through deploying to Leapcell
// Leapcell uses GitHub-based deployments via their web dashboard
func runLeapcellDeploy(ctx context.Context, reader *bufio.Reader, dryRun bool) error {
	if dryRun {
		fmt.Println("\n🔍 DRY RUN: Leapcell Deployment Preview")
		fmt.Println("══════════════════════════════════════")
		fmt.Println("✓ Guide user through Leapcell dashboard setup")
		fmt.Println("✓ Configure GitHub integration")
		fmt.Println("✓ Set build and start commands for Go app")
		fmt.Println("✓ Track deployment URL in .env")
		return nil
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║          🚀 GOTHIC FORGE - LEAPCELL DEPLOYMENT             ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	
	printLeapcellIntro()
	
	// Step-by-step guided setup
	if err := printLeapcellSteps(reader); err != nil {
		return fmt.Errorf("leapcell deployment failed: %w", err)
	}

	return nil
}

func printLeapcellIntro() {
	fmt.Println("\n📘 ABOUT LEAPCELL")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("Leapcell is a modern PaaS platform with:")
	fmt.Println("  • 🎁 20 FREE projects on Hobby tier")
	fmt.Println("  • 🗄️  1 FREE PostgreSQL database")
	fmt.Println("  • ⚡ Serverless-first (pay-per-use)")
	fmt.Println("  • 🔄 Auto-deploy on git push")
	fmt.Println("  • 🌍 Global CDN included")
	fmt.Println("")
	fmt.Println("Perfect for side projects and experimentation!")
	fmt.Println("")
	fmt.Println("⚠️  NOTE: Leapcell deployment is done via their web dashboard.")
	fmt.Println("   This wizard will guide you through the process step-by-step.")
}

func printLeapcellSteps(reader *bufio.Reader) error {
	// Step 1: Prerequisites check
	if err := leapcellStep1Prerequisites(reader); err != nil {
		return err
	}

	// Step 2: Connect GitHub
	if err := leapcellStep2GitHub(reader); err != nil {
		return err
	}

	// Step 3: Create Service
	if err := leapcellStep3CreateService(reader); err != nil {
		return err
	}

	// Step 4: Configure Build Settings
	if err := leapcellStep4Configure(reader); err != nil {
		return err
	}

	// Step 5: Track Deployment URL
	if err := leapcellStep5TrackURL(reader); err != nil {
		return err
	}

	// Success summary
	printLeapcellSuccess()

	return nil
}

func leapcellStep1Prerequisites(reader *bufio.Reader) error {
	fmt.Println("\n══════════════════════════════════════════")
	fmt.Println("STEP 1: Prerequisites")
	fmt.Println("══════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("Before deploying to Leapcell, ensure you have:")
	fmt.Println("")
	fmt.Println("  ✓ Leapcell account (FREE): https://leapcell.io/signup")
	fmt.Println("  ✓ GitHub repository with your Gothic Forge project")
	fmt.Println("  ✓ Git installed and repository pushed to GitHub")
	fmt.Println("")

	// Detect Git remote
	remote := detectGitRemote()
	if remote == "" {
		fmt.Println("⚠️  Warning: Could not detect Git remote")
		fmt.Println("   Make sure your project is in a Git repository and pushed to GitHub")
	} else {
		fmt.Printf("✓ Detected GitHub repository: %s\n", remote)
	}

	fmt.Println("")
	fmt.Print("Have you created a Leapcell account? (y/n): ")
	answer, _ := reader.ReadString('\n')
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(answer)), "y") {
		fmt.Println("")
		fmt.Println("👉 Please create an account at: https://leapcell.io/signup")
		fmt.Println("   Then run this command again.")
		return fmt.Errorf("leapcell account required")
	}

	return nil
}

func leapcellStep2GitHub(reader *bufio.Reader) error {
	fmt.Println("\n══════════════════════════════════════════")
	fmt.Println("STEP 2: Connect GitHub to Leapcell")
	fmt.Println("══════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("1. Go to: https://leapcell.io/new/service")
	fmt.Println("")
	fmt.Println("2. Click 'Connect to GitHub' button")
	fmt.Println("")
	fmt.Println("3. Authorize Leapcell to access your repositories")
	fmt.Println("   - Choose 'All repositories' (recommended)")
	fmt.Println("   - OR 'Only select repositories' (select your Gothic Forge repo)")
	fmt.Println("")
	fmt.Println("4. Complete the GitHub OAuth authorization")
	fmt.Println("")
	fmt.Println("📖 Detailed docs: https://docs.leapcell.io/service/connect-to-github/")
	fmt.Println("")

	fmt.Print("Have you connected your GitHub account? (y/n): ")
	answer, _ := reader.ReadString('\n')
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(answer)), "y") {
		return fmt.Errorf("github connection required")
	}

	return nil
}

func leapcellStep3CreateService(reader *bufio.Reader) error {
	fmt.Println("\n══════════════════════════════════════════")
	fmt.Println("STEP 3: Create Service in Leapcell")
	fmt.Println("══════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("1. Still on https://leapcell.io/new/service")
	fmt.Println("")
	fmt.Println("2. Select your Gothic Forge repository from the list")
	fmt.Println("")
	fmt.Println("3. Choose the branch to deploy (usually 'main' or 'master')")
	fmt.Println("")
	fmt.Println("4. Click 'Continue' or 'Next'")
	fmt.Println("")

	fmt.Print("Press ENTER when you've selected your repository... ")
	reader.ReadString('\n')

	return nil
}

func leapcellStep4Configure(reader *bufio.Reader) error {
	fmt.Println("\n══════════════════════════════════════════")
	fmt.Println("STEP 4: Configure Build Settings")
	fmt.Println("══════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("⚙️  IMPORTANT: Use these EXACT values for Gothic Forge")
	fmt.Println("")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("Runtime:")
	fmt.Println("  Select: Go")
	fmt.Println("")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("Build Command:")
	fmt.Println("  go build -o server ./cmd/server")
	fmt.Println("")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("Start Command:")
	fmt.Println("  ./server")
	fmt.Println("")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("Port:")
	fmt.Println("  8080")
	fmt.Println("")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("")
	fmt.Println("💡 TIP: Leapcell may auto-detect some settings.")
	fmt.Println("   Verify they match the values above!")
	fmt.Println("")

	fmt.Print("Press ENTER when you've configured build settings... ")
	reader.ReadString('\n')

	return nil
}

func leapcellStep5TrackURL(reader *bufio.Reader) error {
	fmt.Println("\n══════════════════════════════════════════")
	fmt.Println("STEP 5: Save Deployment URL")
	fmt.Println("══════════════════════════════════════════")
	fmt.Println("")
	fmt.Println("After deployment completes, you'll receive a URL like:")
	fmt.Println("  https://your-app.leapcell.dev")
	fmt.Println("")
	fmt.Println("Or if you set a custom domain:")
	fmt.Println("  https://yourdomain.com")
	fmt.Println("")
	
	fmt.Print("Enter your Leapcell deployment URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	if url == "" {
		fmt.Println("")
		fmt.Println("⚠️  No URL provided. You can add it to .env manually later:")
		fmt.Println("   LEAPCELL_APP_URL=https://your-app.leapcell.dev")
		return nil
	}

	// Normalize URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Save to .env
	kv := map[string]string{"LEAPCELL_APP_URL": url}
	if err := updateEnvFileInPlace(".env", kv); err != nil {
		fmt.Printf("⚠️  Could not update .env: %v\n", err)
		fmt.Printf("Please add manually: LEAPCELL_APP_URL=%s\n", url)
	} else {
		fmt.Printf("✓ Saved to .env: LEAPCELL_APP_URL=%s\n", url)
	}

	return nil
}

func printLeapcellSuccess() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                   🎉 DEPLOYMENT COMPLETE!                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println("")
	fmt.Println("✅ Your Gothic Forge app is now deployed on Leapcell!")
	fmt.Println("")
	fmt.Println("📋 NEXT STEPS:")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("1. Visit your deployment URL to verify it's working")
	fmt.Println("2. Check the Leapcell dashboard for logs and metrics")
	fmt.Println("3. Set up environment variables in the Leapcell dashboard")
	fmt.Println("   - Go to: Service Settings → Environment Variables")
	fmt.Println("   - Add: DATABASE_URL, JWT_SECRET, etc.")
	fmt.Println("")
	fmt.Println("🔄 AUTO-DEPLOYMENT:")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("Future deployments are automatic!")
	fmt.Println("  git push origin main  ← Triggers automatic deployment")
	fmt.Println("")
	fmt.Println("📊 MONITORING:")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("  • Logs: https://leapcell.io/dashboard")
	fmt.Println("  • Metrics: Built-in analytics and monitoring")
	fmt.Println("  • Alerts: Configure email notifications in settings")
	fmt.Println("")
	fmt.Println("📚 LEARN MORE:")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("  • Docs: https://docs.leapcell.io/")
	fmt.Println("  • Discord: https://discord.gg/qF7efny8x2")
	fmt.Println("  • Examples: https://docs.leapcell.io/examples")
	fmt.Println("")
	fmt.Println("🎁 FREE TIER REMINDER:")
	fmt.Println("──────────────────────────────────────────")
	fmt.Println("You can deploy up to 20 projects for free!")
	fmt.Println("Plus 1 free PostgreSQL database included.")
	fmt.Println("")
}
