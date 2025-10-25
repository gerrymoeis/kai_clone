package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gothicforge3/internal/execx"
	"github.com/spf13/cobra"
)

var (
	workersDeployRun bool
	workersProject   string
)

var deployWorkersCmd = &cobra.Command{
	Use:   "workers",
	Short: "Deploy Cloudflare Workers (wrangler)",
	Long: `Deploy Cloudflare Workers to handle dynamic endpoints.
	
This command deploys the Workers in the workers/ directory using wrangler.
The Workers handle dynamic endpoints (like /counter/sync) that run alongside
Cloudflare Pages (static assets).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		banner()
		
		// Check if workers/ directory exists
		workersDir := "workers"
		if _, err := os.Stat(workersDir); os.IsNotExist(err) {
			fmt.Println("❌ No workers/ directory found")
			fmt.Println("")
			fmt.Println("The workers/ directory should contain:")
			fmt.Println("  • counter.js - Worker handling /counter/sync")
			fmt.Println("  • wrangler.toml - Worker configuration")
			fmt.Println("")
			fmt.Println("Create these files or use 'gforge gen-workers' (coming soon)")
			return nil
		}
		
		// Check if wrangler.toml exists
		wranglerToml := filepath.Join(workersDir, "wrangler.toml")
		if _, err := os.Stat(wranglerToml); os.IsNotExist(err) {
			fmt.Println("❌ No wrangler.toml found in workers/")
			fmt.Println("")
			fmt.Println("Create workers/wrangler.toml with your Worker configuration")
			return nil
		}
		
		// Check for wrangler CLI
		wrPath, wrOK := execx.Look("wrangler")
		if !wrOK {
			fmt.Println("❌ wrangler CLI not found")
			fmt.Println("")
			printWranglerInstallHelp()
			fmt.Println("")
			fmt.Println("Then run:")
			fmt.Println("  gforge deploy workers --run")
			return nil
		}
		
		fmt.Println("✅ Found wrangler:", wrPath)
		fmt.Println("✅ Found workers/wrangler.toml")
		fmt.Println("")
		
		// Deploy command
		ctx := context.Background()
		deployArgs := []string{"deploy"}
		
		// Add project name if specified
		if strings.TrimSpace(workersProject) != "" {
			deployArgs = append(deployArgs, "--name", workersProject)
		}
		
		if workersDeployRun {
			fmt.Println("🚀 Deploying Workers to Cloudflare...")
			fmt.Println("")
			
			// Change to workers directory for deployment
			originalDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			
			if err := os.Chdir(workersDir); err != nil {
				return fmt.Errorf("failed to change to workers/ directory: %w", err)
			}
			defer os.Chdir(originalDir)
			
			// Run wrangler deploy
			cmdLine := "wrangler " + strings.Join(deployArgs, " ")
			fmt.Println("Running:", cmdLine)
			fmt.Println("")
			
			// Build full args for execx.RunInteractive
			fullArgs := append([]string{"wrangler"}, deployArgs...)
			if err := execx.RunInteractive(ctx, cmdLine, fullArgs...); err != nil {
				return fmt.Errorf("wrangler deploy failed: %w", err)
			}
			
			fmt.Println("")
			fmt.Println("✨ Workers deployed successfully!")
			fmt.Println("")
			fmt.Println("Your Workers are now handling:")
			fmt.Println("  • POST /counter/sync - HTMX counter endpoint")
			fmt.Println("")
			fmt.Println("Next steps:")
			fmt.Println("  1. Test your deployed site")
			fmt.Println("  2. Check Worker logs: wrangler tail")
			fmt.Println("  3. Deploy Pages: gforge deploy pages --project=your-project --run")
			
		} else {
			fmt.Println("🔍 Dry-run mode (add --run to deploy)")
			fmt.Println("")
			fmt.Println("Would run:")
			fmt.Println("  cd workers/")
			fmt.Println("  wrangler", strings.Join(deployArgs, " "))
			fmt.Println("")
			fmt.Println("To deploy:")
			fmt.Println("  gforge deploy workers --run")
		}
		
		return nil
	},
}

func init() {
	deployWorkersCmd.Flags().BoolVar(&workersDeployRun, "run", false, "execute wrangler deploy (otherwise dry-run)")
	deployWorkersCmd.Flags().StringVar(&workersProject, "name", "", "Worker name (optional, uses wrangler.toml if not specified)")
	deployCmd.AddCommand(deployWorkersCmd)
}
