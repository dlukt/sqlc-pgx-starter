// Package cli wires the cobra command tree for the application binary.
package cli

import (
	"fmt"
	"os"

	"github.com/dlukt/sqlc-pgx-starter/internal/config"
	"github.com/spf13/cobra"
)

var (
	envPath string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:           "app",
	Short:         "sqlc + pgx starter application",
	Long:          "A starter application demonstrating sqlc, pgx and goose with a cobra CLI.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.Load(envPath)
		if err != nil {
			return err
		}
		cfg = c
		return nil
	},
}

// Execute runs the root command. It is the only entrypoint from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&envPath, "env", ".env", "path to an optional .env file")
}
