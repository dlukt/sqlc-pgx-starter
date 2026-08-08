package cli

import (
	"fmt"

	"github.com/dlukt/sqlc-pgx-starter/internal/database"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migrations (powered by goose)",
}

var migrateUpCmd = &cobra.Command{
	Use:           "up",
	Short:         "Apply all pending migrations",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.Migrate(cfg.DatabaseURL)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the most recent migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.MigrateDown(cfg.DatabaseURL)
	},
}

var migrateRedoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Roll back and re-apply the most recent migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.MigrateRedo(cfg.DatabaseURL)
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the current migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.MigrateStatus(cfg.DatabaseURL)
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration in internal/migrations/sql",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := database.CreateMigration(args[0]); err != nil {
			return err
		}
		fmt.Println("created migration:", args[0])
		return nil
	},
}

func init() {
	migrateCmd.AddCommand(
		migrateUpCmd,
		migrateDownCmd,
		migrateRedoCmd,
		migrateStatusCmd,
		migrateCreateCmd,
	)
	rootCmd.AddCommand(migrateCmd)
}
