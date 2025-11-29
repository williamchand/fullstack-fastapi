package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	genprotov1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (a *App) runHTTP(ctx context.Context) error {
	mux := runtime.NewServeMux()

	// Register handlers for gRPC services
	err := genprotov1.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf(":%s", a.cfg.GRPCPort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return err
	}

	err = genprotov1.RegisterOAuthServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf(":%s", a.cfg.GRPCPort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return err
	}

	err = genprotov1.RegisterDataSourceServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf(":%s", a.cfg.GRPCPort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return err
	}

	err = genprotov1.RegisterBillingServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf(":%s", a.cfg.GRPCPort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return err
	}

	err = genprotov1.RegisterWeddingServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf(":%s", a.cfg.GRPCPort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return err
	}

	handler := a.middleware.Auth.HTTPMiddleware(mux)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.cfg.HTTPPort),
		Handler: handler,
	}

	// Serve
	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	return srv.ListenAndServe()
}
