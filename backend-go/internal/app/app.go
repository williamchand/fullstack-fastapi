package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
)

type App struct {
	cfg           *config.Config
	dbPool        repositories.ConnectionPool
	repos         *Repositories
	services      *AppServices
	middleware    *Middleware
	serviceServer *ServiceServer
}

func NewApp(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	repos, dbPool, err := initRepositories(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	services := initServices(repos)
	middleware, err := initMiddleware(cfg, repos)
	if err != nil {
		return nil, err
	}
	serviceServer := initServiceServer(services)

	return &App{
		cfg:           cfg,
		dbPool:        dbPool,
		repos:         repos,
		services:      services,
		middleware:    middleware,
		serviceServer: serviceServer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.dbPool.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error { return a.runGRPC(ctx) })
	g.Go(func() error { return a.runHTTP(ctx) })
	g.Go(func() error { return a.handleShutdown(ctx, cancel) })

	return g.Wait()
}

func (a *App) handleShutdown(ctx context.Context, cancel context.CancelFunc) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case sig := <-c:
		log.Println("Received:", sig)
		cancel()
		time.Sleep(3 * time.Second)
		return nil
	}
}
