# sqlc-pgx-starter

A production-oriented starter template for a Go service backed by PostgreSQL,
using **[sqlc]** for type-safe SQL, **[pgx]** for the driver/connection pool,
**[goose]** for migrations (embedded in the binary), and **[Cobra]** for the CLI.

Clone it, run `./rename.sh`, and start writing queries.

[sqlc]: https://docs.sqlc.dev
[pgx]: https://github.com/jackc/pgx
[goose]: https://github.com/pressly/goose
[Cobra]: https://github.com/spf13/cobra

## Features

- **sqlc** — write SQL, get fully-typed Go (pgx v5 backend). No ORM.
- **pgx** — connection pool with sane defaults (`internal/database`).
- **goose** — versioned migrations embedded via `//go:embed`, applied by the CLI.
- **Cobra CLI** — one binary with `serve`, `migrate ...`, and `version` subcommands.
- **testcontainers-go** — an end-to-end integration test that spins up a real Postgres.
- **docker-compose** — one-command local Postgres.
- **Makefile** — every common task wrapped up.
- **rename.sh** — rebrand the whole repo to your module path in one shot.
- **GitHub Actions** — CI that builds, vets, checks generated-code drift, and tests.

## Prerequisites

| Tool     | Why                          | Arch               | Debian / Ubuntu         | Fedora / RHEL        |
|----------|------------------------------|--------------------|-------------------------|----------------------|
| Go 1.26+ | language                     | `pacman -S go`     | `apt install golang-go` | `dnf install golang` |
| Docker   | local DB + integration tests | `pacman -S docker` | `apt install docker.io` | `dnf install docker` |
| sqlc     | code generation              | `pacman -S sqlc`   | `go install …@latest` ¹ | `go install …@latest` ¹ |

¹ `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` — sqlc isn't in the
Debian/Fedora package repos, so install it through Go.

sqlc is only needed to regenerate code; the generated output is committed, so
building/testing works without it. If your distro ships an older Go, install
1.26+ from the [official binaries](https://go.dev/dl/) instead.

## Quickstart

One-liner (assumes Docker is running):

```bash
./rename.sh github.com/yourorg/your-service && \
  cp .env.example .env && make db-up && make migrate-up && make run ARGS=serve
```

Step by step:

```bash
# 1. Rename the repo to your own project.
./rename.sh github.com/yourorg/your-service

# 2. Start a local Postgres.
cp .env.example .env
make db-up

# 3. Apply migrations.
make migrate-up

# 4. Run the example command (migrates, connects, inserts a row, prints count).
make run ARGS=serve   # or: go run ./cmd/app serve

# 5. Run the tests (requires Docker).
make test
```

## CLI

The binary is built from `cmd/app` and exposes a Cobra command tree:

```text
app serve              # run the app (migrate -> connect -> example query)
app migrate up         # apply all pending migrations
app migrate down       # roll back the most recent migration
app migrate redo       # roll back + re-apply the most recent migration
app migrate status     # print migration status
app migrate create NAME  # scaffold a new migration in internal/migrations/sql
app version            # print version (set via -ldflags at build time)
```

Run via `go run ./cmd/app <args>` or `make run ARGS="<args>"`. Build a versioned binary:

```bash
make build VERSION=v0.1.0
./bin/app version   # v0.1.0
```

## Configuration

Configuration is read from a `.env` file (optional) and the process environment.
See `.env.example`. Either set a full DSN:

```dotenv
DATABASE_URL=postgres://user:pass@localhost:5432/appdb?sslmode=disable
```

…or leave `DATABASE_URL` empty and set the individual fields, which are composed
into a DSN for you: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`,
`DB_SSLMODE`, plus `APP_ENV` and `HTTP_ADDR`.

## The sqlc workflow

sqlc reads SQL and emits type-safe Go. The loop is:

1. **Edit schema** in a goose migration under `internal/migrations/sql/`.
2. **Edit queries** in `db/queries/*.sql` using the
   `-- name: FunctionName :one|:many|:exec` annotation.
3. **Regenerate:** `make generate` (writes to `internal/db/`).
4. **Use** the generated `db.New(pool)` in your code.

`sqlc.yaml` is configured for the `pgx/v5` backend, emits an interface, JSON
tags, pointers for nullable columns, and maps `uuid` columns to
`github.com/google/uuid.UUID`.

```go
pool, _ := database.New(ctx, cfg.DatabaseURL)
q := db.New(pool)

user, err := q.CreateUser(ctx, db.CreateUserParams{
    Name:  "Ada Lovelace",
    Email: "ada@example.com",
})
```

To confirm your generated code matches the SQL without a database, run
`make sqlc-verify` (regenerates and fails on any drift).

## Migrations

Migrations are goose-formatted SQL files in `internal/migrations/sql/` and are
**embedded into the binary** via `internal/migrations/embed.go`, so the
production binary is fully self-contained — no `.sql` files shipped alongside.

The `database` package wraps goose:

| Function               | CLI                       |
|------------------------|---------------------------|
| `database.Migrate`     | `app migrate up`          |
| `database.MigrateDown` | `app migrate down`        |
| `database.MigrateRedo` | `app migrate redo`        |
| `database.MigrateStatus` | `app migrate status`    |
| `database.CreateMigration` | `app migrate create NAME` |

Create a new migration (writes to the source tree):

```bash
make new-migration NAME=add_users_index
```

## Tests

`internal/database/database_test.go` is a full end-to-end test: it launches an
ephemeral Postgres via **testcontainers-go**, runs the goose migrations, and
exercises every generated query through a pgx pool.

```bash
make test         # full suite (needs Docker daemon running)
make test-short   # skip the Docker-based integration test
```

## Project layout

```
.
├── cmd/app/main.go              # entrypoint: builds the Cobra command tree
├── db/queries/*.sql             # sqlc query input
├── internal/
│   ├── cli/                     # Cobra commands (root, serve, migrate, version)
│   ├── config/                  # env/.env configuration
│   ├── database/                # pgx pool + goose migration helpers
│   ├── db/                      # generated by sqlc (do not edit)
│   └── migrations/
│       ├── embed.go             # //go:embed of sql/
│       └── sql/*.sql            # goose migrations (source of truth)
├── docker-compose.yml           # local Postgres
├── sqlc.yaml                    # sqlc configuration
├── Makefile                     # task runner
├── rename.sh                    # rebrand the repo after cloning
└── .github/workflows/ci.yml     # CI
```

## Customizing

- **Module path / repo name:** `./rename.sh github.com/yourorg/your-service`
- **Connection pool sizing:** edit defaults in `internal/database/database.go`.
- **sqlc options** (null pointers, interfaces, type overrides): edit `sqlc.yaml`.
- **Version:** injected at build via `-ldflags "-X ...cli.Version=..."` (see `Makefile`).

## License

This starter is provided as-is for you to adapt; add your own `LICENSE` file.
