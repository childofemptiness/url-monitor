package app

import (
	"context"
	"fmt"
	"net/http"
	apphttp "url-monitor/internal/http"
	"url-monitor/internal/monitor"
	"url-monitor/internal/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	server *http.Server
	db     *pgxpool.Pool
}

func New(ctx context.Context, addr, databaseURL string) (*App, error) {
	pool, err := postgres.NewPool(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	monitorRepo    := postgres.NewMonitorRepository(pool)
	monitorService := monitor.NewMonitorService(monitorRepo)
	handler        := apphttp.NewHandler(monitorService)
	router         := apphttp.NewRouter(handler)

	server := &http.Server{
		Addr: addr,
		Handler: router,
	}

	return &App{
		server: server,
		db: pool,
	}, nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Close(ctx context.Context) error {
	var result error

	if err := a.server.Shutdown(ctx); err != nil {
		result = fmt.Errorf("shutdown server: %w", err)
	}

	a.db.Close()

	return result
}
