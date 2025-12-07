package app

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

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

	err = genprotov1.RegisterPublicServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf(":%s", a.cfg.GRPCPort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return err
	}

	handler := a.middleware.Auth.HTTPMiddleware(mux)

	// Root mux to serve OpenAPI specs without auth and gRPC-Gateway with auth
	rootMux := http.NewServeMux()

	// Serve static OpenAPI swagger JSON files from gen/openapi/v1 at /openapi/v1/
	openapiDir := filepath.Join("gen", "openapi", "v1")
	fs := http.FileServer(http.Dir(openapiDir))
	rootMux.Handle("/v1/openapi/", http.StripPrefix("/v1/openapi/", fs))

	// All other routes go through auth + grpc-gateway
	rootMux.Handle("/", handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.cfg.HTTPPort),
		Handler: rootMux,
	}

	// Serve
	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	return srv.ListenAndServe()
}
