package cli

import (
	"context"
	"fmt"

	"github.com/dlukt/sqlc-pgx-starter/internal/database"
	"github.com/dlukt/sqlc-pgx-starter/internal/db"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the application",
	Long:  "Runs migrations, opens a database pool and demonstrates a sqlc query end to end.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if err := database.Migrate(cfg.DatabaseURL); err != nil {
			return err
		}

		pool, err := database.New(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()

		q := db.New(pool)

		u, err := q.CreateUser(ctx, db.CreateUserParams{
			Name:  "Ada Lovelace",
			Email: "ada@example.com",
		})
		if err != nil {
			return err
		}
		fmt.Printf("created user %s <%s>\n", u.Name, u.Email)

		count, err := q.CountUsers(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("total users: %d\n", count)
		fmt.Printf("listening on %s (stub server)\n", cfg.HTTPAddr)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
