package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"

	"github.com/wreckitral/production-backend-go/internal/platform/config"
	"github.com/wreckitral/production-backend-go/internal/platform/server"
)

func newTestServer(t *testing.T) *httptest.Server {
    t.Helper()

    cfg := config.Config{
        HTTP: config.HTTP{
            ReadHeaderTimeout: 5 * time.Second,
            ReadTimeout:       5 * time.Second,
            WriteTimeout:      10 * time.Second,
            IdleTimeout:       120 * time.Second,
            ShutdownTimeout:   1 * time.Second,
            MaxHeaderBytes:    1 << 20,
        },
        JWT: config.JWT{
            Secret: "dev-only-secret-change-me-in-production-please-32b",
            TTL:    time.Hour,
            Issuer: "blog",
        },
    }

    ln, err := net.Listen("tcp", ":0")
    if err != nil { t.Fatalf("listen: %v", err) }
    t.Cleanup(func() { _ = ln.Close() })

    app := server.New(server.Deps{
        Config:   cfg,
        Logger:   testLogger(),
        DB:       testPool,
        Listener: ln,
        WebFS:    fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte("ok")}},
    })

    ts := httptest.NewUnstartedServer(app.Handler())
    ts.Start()
    return ts
}

func TestPostsAPI_CreateAndList(t *testing.T) {
    cleanDB(t)
    ts := newTestServer(t)
    defer ts.Close()

    body := mustJSON(t, map[string]string{"email": "alice@example.com", "password": "supersecret"})
    res, err := http.Post(ts.URL+"/api/auth/register", "application/json", bytes.NewReader(body))
    if err != nil { t.Fatalf("register request: %v", err) }
    defer res.Body.Close()
    if res.StatusCode != http.StatusCreated {
        t.Fatalf("register: got %d, body=%s", res.StatusCode, readAll(res.Body))
    }

    body = mustJSON(t, map[string]string{"email": "alice@example.com", "password": "supersecret"})
    res, err = http.Post(ts.URL+"/api/auth/login", "application/json", bytes.NewReader(body))
    if err != nil { t.Fatalf("login request: %v", err) }
    defer res.Body.Close()
    if res.StatusCode != http.StatusOK {
        t.Fatalf("login: got %d, body=%s", res.StatusCode, readAll(res.Body))
    }
    var login struct{ Token string `json:"token"` }
    decode(t, res.Body, &login)

    body = mustJSON(t, map[string]string{"title": "hello", "body": "world"})
    req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, ts.URL+"/api/posts", bytes.NewReader(body))
    if err != nil { t.Fatalf("new request: %v", err) }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+login.Token)
    res, err = http.DefaultClient.Do(req)
    if err != nil { t.Fatalf("create request: %v", err) }
    defer res.Body.Close()
    if res.StatusCode != http.StatusCreated {
        t.Fatalf("create: got %d, body=%s", res.StatusCode, readAll(res.Body))
    }

    res, err = http.Get(ts.URL + "/api/posts")
    if err != nil { t.Fatalf("list request: %v", err) }
    defer res.Body.Close()
    if res.StatusCode != http.StatusOK {
        t.Fatalf("list: got %d, body=%s", res.StatusCode, readAll(res.Body))
    }
    var listed []map[string]any
    decode(t, res.Body, &listed)
    if len(listed) != 1 || listed[0]["title"] != "hello" {
        t.Fatalf("list: %+v", listed)
    }
}

func mustJSON(t *testing.T, v any) []byte {
    t.Helper()
    b, err := json.Marshal(v)
    if err != nil { t.Fatal(err) }
    return b
}

func decode(t *testing.T, r io.Reader, v any) {
    t.Helper()
    if err := json.NewDecoder(r).Decode(v); err != nil { t.Fatal(err) }
}

func readAll(r io.Reader) string { b, _ := io.ReadAll(r); return string(b) }
