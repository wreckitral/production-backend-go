package tests

import (
    "context"
    "database/sql"
    "fmt"
    "io"
    "log"
    "log/slog"
    "os"
    "testing"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/ory/dockertest/v3"
    "github.com/ory/dockertest/v3/docker"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
    pool, err := dockertest.NewPool("")
    if err != nil { log.Fatalf("dockertest pool: %s", err) }
    if err := pool.Client.Ping(); err != nil { log.Fatalf("docker not reachable: %s", err) }

    resource, err := pool.RunWithOptions(&dockertest.RunOptions{
        Repository: "postgres",
        Tag:        "16-alpine",
        Env: []string{"POSTGRES_USER=blog", "POSTGRES_PASSWORD=blog", "POSTGRES_DB=blog"},
    }, func(c *docker.HostConfig) {
        c.AutoRemove = true
        c.RestartPolicy = docker.RestartPolicy{Name: "no"}
    })
    if err != nil { log.Fatalf("start postgres: %s", err) }

    hostPort := resource.GetPort("5432/tcp")
    dsn := fmt.Sprintf("postgres://blog:blog@localhost:%s/blog?sslmode=disable", hostPort)

    pool.MaxWait = 60 * time.Second
    if err := pool.Retry(func() error {
        db, err := sql.Open("pgx", dsn)
        if err != nil { return err }
        defer db.Close()
        return db.Ping()
    }); err != nil { log.Fatalf("postgres did not come up: %s", err) }

    mig, err := migrate.New("file://../migrations", dsn)
    if err != nil { log.Fatalf("migrate init: %s", err) }
    if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("migrate up: %s", err)
    }

    testPool, err = pgxpool.New(context.Background(), dsn)
    if err != nil { log.Fatalf("pgxpool: %s", err) }

    code := m.Run()
    testPool.Close()
    _ = pool.Purge(resource)
    os.Exit(code)
}

func cleanDB(t *testing.T) {
    t.Helper()
    _, err := testPool.Exec(context.Background(),
        `TRUNCATE users, posts RESTART IDENTITY CASCADE`)
    if err != nil { t.Fatalf("truncate: %v", err) }
}

func seedUser(t *testing.T, email string) uuid.UUID {
    t.Helper()

    var id uuid.UUID
    err := testPool.QueryRow(context.Background(), `
        INSERT INTO users (email, password_hash)
        VALUES (lower($1), 'test-hash')
        RETURNING id
    `, email).Scan(&id)
    if err != nil { t.Fatalf("seed user: %v", err) }
    return id
}

func testLogger() *slog.Logger {
    return slog.New(slog.NewTextHandler(io.Discard, nil))
}
