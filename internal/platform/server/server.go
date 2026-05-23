package server

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wreckitral/production-backend-go/internal/auth"
	"github.com/wreckitral/production-backend-go/internal/middleware"
	"github.com/wreckitral/production-backend-go/internal/platform/config"
	"github.com/wreckitral/production-backend-go/internal/post"
)

type Deps struct {
	Config   config.Config
	Logger   *slog.Logger
	DB       *pgxpool.Pool
	Listener net.Listener
	WebFS    fs.FS
}

type App struct {
	deps   Deps
	server *http.Server
}

func New(d Deps) *App {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recover(d.Logger))
	r.Use(middleware.Logger(d.Logger))

	jwtMW := middleware.NewJWT(d.Config.JWT.Secret, d.Config.JWT.Issuer)

	r.Get("/healthz", healthz)
	r.Get("/readyz", readyz(d.DB))

	authRepo := auth.NewRepo(d.DB)
	authSvc := auth.NewService(authRepo, d.Config.JWT)
	authH := auth.NewHandler(authSvc)
	auth.Route(r, authH)

	postRepo := post.NewRepo(d.DB)
	postSvc := post.NewService(postRepo)
	postH := post.NewHandler(postSvc)
	post.Route(r, postH, jwtMW)

	if d.WebFS != nil {
		web := d.WebFS
		if sub, err := fs.Sub(d.WebFS, "web"); err == nil {
			web = sub
		}
		r.Handle("/*", http.FileServer(http.FS(web)))
	}

	return &App{
		deps: d,
		server: &http.Server{
			Handler:           r,
			ReadHeaderTimeout: d.Config.HTTP.ReadHeaderTimeout,
			ReadTimeout:       d.Config.HTTP.ReadTimeout,
			WriteTimeout:      d.Config.HTTP.WriteTimeout,
			IdleTimeout:       d.Config.HTTP.IdleTimeout,
			MaxHeaderBytes:    d.Config.HTTP.MaxHeaderBytes,
		},
	}
}

func (a *App) Run(ctx context.Context) error {
	serveErr := make(chan error, 1)
	go func() {
		if err := a.server.Serve(a.deps.Listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
			return
		}
		serveErr <- nil
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		a.deps.Logger.Info("shutdown requested")
		shutCtx, cancel := context.WithTimeout(context.Background(), a.deps.Config.HTTP.ShutdownTimeout)
		defer cancel()
		return a.server.Shutdown(shutCtx)
	}
}

func (a *App) Handler() http.Handler { return a.server.Handler }

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func readyz(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"not ready"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	}
}
