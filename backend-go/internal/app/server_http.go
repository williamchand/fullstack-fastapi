package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	genprotov1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func (a *App) runHTTP(ctx context.Context) error {
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(gatewayErrorHandler),
	)

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

	// Apply CORS middleware to all routes
	corsHandler := cors.Middleware(a.cfg.Security.CORSAllowedOrigins)(rootMux)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.cfg.HTTPPort),
		Handler: corsHandler,
	}

	// Serve
	go func() {
		<-ctx.Done()
		// Graceful shutdown
		srv.Shutdown(context.Background())
	}()

	return srv.ListenAndServe()
}

// gatewayErrorHandler ensures gRPC-Gateway sends JSON errors with a consistent shape.
// It includes at least a "message" field so the frontend can display backend-provided messages.
func gatewayErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	st := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(st.Code())

	// Build a simple JSON error payload recognized by the frontend
	type errorResponse struct {
		Code    int         `json:"code,omitempty"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	}

	resp := errorResponse{
		Message: st.Message(),
	}
	if st.Code() != codes.OK {
		resp.Code = int(st.Code())
	}
	if len(st.Details()) > 0 {
		resp.Details = st.Details()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(resp)
}
