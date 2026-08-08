package database_test

import (
	"context"
	"testing"

	"github.com/dlukt/sqlc-pgx-starter/internal/database"
	"github.com/dlukt/sqlc-pgx-starter/internal/db"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

// TestCRUD exercises the full stack end to end against an ephemeral Postgres:
// goose migrations + a pgx pool + sqlc-generated queries. Skipped under -short.
func TestCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}

	ctx := context.Background()

	pgC, err := tcpg.Run(ctx,
		"docker.io/postgres:18-alpine",
		// Wait for Postgres to be fully ready: the init server prints the
		// readiness line once, then the real server prints it again, so we
		// wait for the second occurrence before connecting.
		tcpg.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := pgC.Terminate(ctx); err != nil {
			t.Logf("terminate container: %v", err)
		}
	})

	dsn, err := pgC.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	if err := database.Migrate(dsn); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	pool, err := database.New(ctx, dsn)
	if err != nil {
		t.Fatalf("open pool: %v", err)
	}
	defer pool.Close()

	q := db.New(pool)

	// Create (nullable bio left nil).
	u, err := q.CreateUser(ctx, db.CreateUserParams{
		Name:  "Ada",
		Email: "ada@example.com",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if u.Bio != nil {
		t.Fatalf("bio = %v, want nil", *u.Bio)
	}

	// Read by id.
	got, err := q.GetUser(ctx, u.ID)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if got.Name != "Ada" {
		t.Fatalf("name = %q, want %q", got.Name, "Ada")
	}

	// List + count.
	if users, err := q.ListUsers(ctx, 10); err != nil {
		t.Fatalf("list users: %v", err)
	} else if len(users) != 1 {
		t.Fatalf("len(users) = %d, want 1", len(users))
	}
	if n, err := q.CountUsers(ctx); err != nil {
		t.Fatalf("count users: %v", err)
	} else if n != 1 {
		t.Fatalf("count = %d, want 1", n)
	}

	// Update (set nullable bio).
	bio := "mathematician"
	updated, err := q.UpdateUser(ctx, db.UpdateUserParams{
		ID:   u.ID,
		Name: "Ada L.",
		Bio:  &bio,
	})
	if err != nil {
		t.Fatalf("update user: %v", err)
	}
	if updated.Bio == nil || *updated.Bio != bio {
		t.Fatalf("bio = %v, want %q", updated.Bio, bio)
	}

	// Delete + verify count drops to zero.
	if err := q.DeleteUser(ctx, u.ID); err != nil {
		t.Fatalf("delete user: %v", err)
	}
	if n, err := q.CountUsers(ctx); err != nil {
		t.Fatalf("count users after delete: %v", err)
	} else if n != 0 {
		t.Fatalf("count = %d, want 0", n)
	}
}
