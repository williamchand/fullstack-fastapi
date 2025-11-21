package app

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/delivery/grpc"
)

func (a *App) runHTTP(ctx context.Context) error {
	return grpc.RunHTTPServer(ctx, a.cfg.HTTPPort, a.cfg.GRPCPort, a.middleware.Auth)
}
