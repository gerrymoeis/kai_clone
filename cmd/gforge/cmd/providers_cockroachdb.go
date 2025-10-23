package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// CockroachDB Cloud API Documentation:
// https://www.cockroachlabs.com/docs/cockroachcloud/cloud-api
//
// Why CockroachDB Serverless?
// 1. PostgreSQL-compatible (uses pgx driver, no code changes)
// 2. True serverless pricing (pay only for what you use)
// 3. Global distribution built-in (low latency worldwide)
// 4. Automatic scaling (0 to massive without configuration)
// 5. Built-in backups and point-in-time recovery
// 6. ACID compliance with distributed architecture
//
// Educational Value:
// Users learn about distributed SQL databases and modern cloud-native architecture

const cockroachAPIBase = "https://cockroachlabs.cloud/api/v1"

type cockroachClient struct {
	apiKey string
	hc     *http.Client
}

func newCockroachClientFromEnv() (*cockroachClient, error) {
	apiKey := strings.TrimSpace(os.Getenv("COCKROACH_API_KEY"))
	if apiKey == "" {
		return nil, errors.New("COCKROACH_API_KEY is not set")
	}
	return &cockroachClient{
		apiKey: apiKey,
		hc:     &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *cockroachClient) do(ctx context.Context, method, path string, in any, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, cockroachAPIBase+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cockroachdb api %s %s: %s: %s", method, path, resp.Status, string(b))
	}
	if out != nil && resp.StatusCode != 204 {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

// CockroachDB API Types
type cockroachCluster struct {
	ID              string                  `json:"id"`
	Name            string                  `json:"name"`
	CockroachVersion string                  `json:"cockroach_version"`
	Plan            string                  `json:"plan"`
	CloudProvider   string                  `json:"cloud_provider"`
	State           string                  `json:"state"`
	Config          cockroachClusterConfig  `json:"config"`
	Regions         []cockroachClusterRegion `json:"regions"`
}

type cockroachClusterConfig struct {
	Serverless cockroachServerlessConfig `json:"serverless"`
}

type cockroachServerlessConfig struct {
	SpendLimit int    `json:"spend_limit"`
	RoutingID  string `json:"routing_id"`
}

type cockroachClusterRegion struct {
	Name      string `json:"name"`
	SQLDns    string `json:"sql_dns"`
	UIDns     string `json:"ui_dns"`
}

// cockroachInteractiveProvision guides user through CockroachDB Serverless setup.
// This is the opinionated Gothic Forge database provider.
func cockroachInteractiveProvision(ctx context.Context, dryRun bool) (string, error) {
	// Check if already provisioned
	cur := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if cur != "" {
		fmt.Println("  • CockroachDB: DATABASE_URL already set")
		return cur, nil
	}

	if dryRun {
		fmt.Println("  • CockroachDB (dry-run): would provision serverless cluster and save DATABASE_URL")
		return "", nil
	}

	// Check for API key
	apiKey := strings.TrimSpace(os.Getenv("COCKROACH_API_KEY"))
	if apiKey == "" {
		fmt.Println("  • CockroachDB: COCKROACH_API_KEY not set")
		fmt.Println("")
		fmt.Println("    How to get your API key:")
		fmt.Println("    1. Sign up: https://cockroachlabs.cloud/signup")
		fmt.Println("    2. Create API key: https://cockroachlabs.cloud/account/api-access")
		fmt.Println("    3. Set in .env: COCKROACH_API_KEY=<your-key>")
		fmt.Println("")
		fmt.Println("    Educational Note:")
		fmt.Println("    CockroachDB Serverless is Gothic Forge's opinionated database choice because:")
		fmt.Println("    - PostgreSQL-compatible (no code changes from traditional Postgres)")
		fmt.Println("    - True serverless (pay only for usage, scales to zero)")
		fmt.Println("    - Global distribution (low latency worldwide)")
		fmt.Println("    - Built-in resilience (automatic replication and failover)")
		fmt.Println("")
		return "", errors.New("COCKROACH_API_KEY required")
	}

	client, err := newCockroachClientFromEnv()
	if err != nil {
		return "", err
	}

	fmt.Println("  • CockroachDB: Provisioning serverless cluster...")
	fmt.Println("    → This will create:")
	fmt.Println("      - Serverless cluster (auto-scaling)")
	fmt.Println("      - Database and SQL user")
	fmt.Println("      - Secure connection string")

	reader := bufio.NewReader(os.Stdin)

	// Get cluster name
	fmt.Print("    Cluster name [gothic-forge-db]: ")
	clusterName, _ := reader.ReadString('\n')
	clusterName = strings.TrimSpace(clusterName)
	if clusterName == "" {
		clusterName = "gothic-forge-db"
	}

	// Get region
	fmt.Println("    Available regions:")
	fmt.Println("      - us-east-1 (N. Virginia, AWS)")
	fmt.Println("      - us-west-2 (Oregon, AWS)")
	fmt.Println("      - eu-central-1 (Frankfurt, AWS)")
	fmt.Println("      - ap-southeast-1 (Singapore, AWS)")
	fmt.Print("    Region [us-east-1]: ")
	region, _ := reader.ReadString('\n')
	region = strings.TrimSpace(region)
	if region == "" {
		region = "us-east-1"
	}

	// Create cluster
	fmt.Println("    → Creating serverless cluster (this may take 30-60 seconds)...")
	cluster, err := createCockroachServerlessCluster(ctx, client, clusterName, region)
	if err != nil {
		return "", fmt.Errorf("failed to create cluster: %w", err)
	}
	fmt.Printf("    → Cluster created: %s (ID: %s)\n", cluster.Name, cluster.ID)

	// Wait for cluster to be ready
	fmt.Println("    → Waiting for cluster to be ready...")
	if err := waitForClusterReady(ctx, client, cluster.ID); err != nil {
		return "", fmt.Errorf("cluster failed to become ready: %w", err)
	}
	fmt.Println("    → Cluster is ready")

	// Create SQL user
	username := "gothicforge"
	fmt.Printf("    → Creating SQL user: %s\n", username)
	password, err := createCockroachSQLUser(ctx, client, cluster.ID, username)
	if err != nil {
		return "", fmt.Errorf("failed to create SQL user: %w", err)
	}

	// Create database
	dbName := "app"
	fmt.Printf("    → Creating database: %s\n", dbName)
	if err := createCockroachDatabase(ctx, client, cluster.ID, dbName); err != nil {
		return "", fmt.Errorf("failed to create database: %w", err)
	}

	// Build connection string
	host := cluster.Regions[0].SQLDns
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:26257/%s?sslmode=verify-full",
		username, password, host, dbName)

	// Save to .env
	kv := map[string]string{"DATABASE_URL": dsn}
	if err := updateEnvFileInPlace(".env", kv); err != nil {
		// Fallback: append
		if f, ferr := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600); ferr == nil {
			defer f.Close()
			_, _ = f.WriteString("\n# CockroachDB Serverless\nDATABASE_URL=" + dsn + "\n")
		}
	}
	_ = os.Setenv("DATABASE_URL", dsn)

	fmt.Println("    → DATABASE_URL saved to .env")
	fmt.Println("")
	fmt.Println("  • CockroachDB provisioning complete!")
	fmt.Println("    → Console: https://cockroachlabs.cloud")
	fmt.Printf("    → Cluster: %s\n", cluster.Name)
	fmt.Println("")

	// Auto-run migrations
	if err := runMigrationsAuto(ctx, dsn); err != nil {
		fmt.Printf("    ⚠️  Migration warning: %v\n", err)
		fmt.Println("    → You can run migrations manually with: gforge db --migrate")
	}

	fmt.Println("")
	fmt.Println("    Educational Note:")
	fmt.Println("    Your data is now distributed across multiple availability zones")
	fmt.Println("    for high availability. CockroachDB automatically handles:")
	fmt.Println("    - Replication (data safety)")
	fmt.Println("    - Load balancing (performance)")
	fmt.Println("    - Failover (resilience)")
	fmt.Println("")

	return dsn, nil
}

// createCockroachServerlessCluster creates a new serverless cluster.
func createCockroachServerlessCluster(ctx context.Context, client *cockroachClient, name, region string) (*cockroachCluster, error) {
	payload := map[string]interface{}{
		"name":           name,
		"provider":       "AWS",
		"spec": map[string]interface{}{
			"serverless": map[string]interface{}{
				"regions": []string{region},
				"spend_limit": 0, // Free tier (can be adjusted)
			},
		},
	}

	var result struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		CloudProvider   string `json:"cloud_provider"`
		CockroachVersion string `json:"cockroach_version"`
		Plan            string `json:"plan"`
		State           string `json:"state"`
		Regions         []cockroachClusterRegion `json:"regions"`
		Config          cockroachClusterConfig `json:"config"`
	}

	if err := client.do(ctx, "POST", "/clusters", payload, &result); err != nil {
		return nil, err
	}

	return &cockroachCluster{
		ID:              result.ID,
		Name:            result.Name,
		CockroachVersion: result.CockroachVersion,
		Plan:            result.Plan,
		CloudProvider:   result.CloudProvider,
		State:           result.State,
		Regions:         result.Regions,
		Config:          result.Config,
	}, nil
}

// waitForClusterReady polls until cluster state is CREATED.
func waitForClusterReady(ctx context.Context, client *cockroachClient, clusterID string) error {
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return errors.New("timeout waiting for cluster to be ready")
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			var cluster cockroachCluster
			if err := client.do(ctx, "GET", "/clusters/"+clusterID, nil, &cluster); err != nil {
				return err
			}
			if cluster.State == "CREATED" {
				return nil
			}
			// States: CREATING, CREATED, CREATION_FAILED
			if strings.Contains(strings.ToUpper(cluster.State), "FAILED") {
				return fmt.Errorf("cluster creation failed: %s", cluster.State)
			}
		}
	}
}

// createCockroachSQLUser creates a SQL user and returns the generated password.
func createCockroachSQLUser(ctx context.Context, client *cockroachClient, clusterID, username string) (string, error) {
	payload := map[string]string{
		"name": username,
	}

	var result struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := client.do(ctx, "POST", "/clusters/"+clusterID+"/sql-users", payload, &result); err != nil {
		return "", err
	}

	return result.Password, nil
}

// createCockroachDatabase creates a database in the cluster.
func createCockroachDatabase(ctx context.Context, client *cockroachClient, clusterID, dbName string) error {
	payload := map[string]string{
		"name": dbName,
	}

	return client.do(ctx, "POST", "/clusters/"+clusterID+"/databases", payload, nil)
}
