package cmd

import (
  "context"
  "database/sql"
  "errors"
  "fmt"
  "os"
  "path/filepath"
  "time"

  "github.com/pressly/goose/v3"
  _ "github.com/jackc/pgx/v5/stdlib"
  "github.com/spf13/cobra"
)

var (
  dbMigrate bool
  dbReset   bool
  dbStatus  bool
)

var dbCmd = &cobra.Command{
  Use:   "db",
  Short: "Database helpers (Neon)",
  RunE: func(cmd *cobra.Command, args []string) error {
    banner()
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
      return errors.New("DATABASE_URL is not set; cannot run migrations")
    }
    dir := filepath.Join("app", "db", "migrations")
    if _, err := os.Stat(dir); os.IsNotExist(err) {
      return fmt.Errorf("migrations directory not found: %s", dir)
    }
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    // Open pgx driver via database/sql
    dbx, err := sql.Open("pgx", dsn)
    if err != nil { return err }
    defer dbx.Close()
    if err := dbx.PingContext(ctx); err != nil { return err }
    // Run action
    if dbStatus {
      fmt.Println("DB: status...")
      if err := goose.Status(dbx, dir); err != nil { return err }
      return nil
    }
    if dbMigrate {
      fmt.Println("DB: applying migrations...")
      if err := goose.Up(dbx, dir); err != nil { return err }
      fmt.Println("DB: migrations complete")
      return nil
    }
    if dbReset {
      fmt.Println("DB: resetting database via goose reset...")
      if err := goose.Reset(dbx, dir); err != nil { return err }
      fmt.Println("DB: reset complete")
      return nil
    }
    fmt.Println("Usage: gforge db --status | --migrate | --reset")
    return nil
  },
}

func init() {
  dbCmd.Flags().BoolVar(&dbMigrate, "migrate", false, "apply migrations")
  dbCmd.Flags().BoolVar(&dbReset, "reset", false, "reset database")
  dbCmd.Flags().BoolVar(&dbStatus, "status", false, "show migration status")
  rootCmd.AddCommand(dbCmd)
}

// runMigrationsAuto automatically applies database migrations after provisioning.
// This is called by deployment workflows to ensure schema is up-to-date.
// Returns nil if migrations succeed or if migrations directory doesn't exist yet.
func runMigrationsAuto(ctx context.Context, dsn string) error {
  if dsn == "" {
    return nil // No database configured, skip migrations
  }

  dir := filepath.Join("app", "db", "migrations")
  if _, err := os.Stat(dir); os.IsNotExist(err) {
    fmt.Println("    → No migrations directory found (this is OK for new projects)")
    return nil // Migrations directory doesn't exist yet, that's fine
  }

  fmt.Println("    → Running database migrations...")

  // Open database connection
  dbx, err := sql.Open("pgx", dsn)
  if err != nil {
    return fmt.Errorf("failed to connect to database: %w", err)
  }
  defer dbx.Close()

  // Test connection
  if err := dbx.PingContext(ctx); err != nil {
    return fmt.Errorf("database unreachable: %w", err)
  }

  // Apply migrations
  if err := goose.Up(dbx, dir); err != nil {
    return fmt.Errorf("migrations failed: %w", err)
  }

  fmt.Println("    → Migrations applied successfully")
  return nil
}
