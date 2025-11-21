package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"myapp/internal/delivery/grpc"
	"myapp/internal/domain/services"
	"myapp/internal/infrastructure/auth"
	"myapp/internal/infrastructure/database"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg := loadConfig()

	// Setup database
	dbPool := setupDatabase(ctx, cfg.DatabaseURL)
	defer dbPool.Close()

	// Initialize repositories
	queries := database.New(dbPool)
	userRepo := database.NewUserRepository(queries)
	oauthRepo := database.NewOAuthRepository(queries)

	// Initialize services
	userService := services.NewUserService(userRepo, oauthRepo)

	// Initialize auth
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	roleValidator := auth.NewRoleValidator()
	authMiddleware := auth.NewAuthMiddleware(jwtService, roleValidator, userRepo)

	// Initialize gRPC servers
	userServer := grpc.NewUserServer(userService)

	// Start servers
	if err := runServers(ctx, cfg, userServer, authMiddleware); err != nil {
		log.Fatal("Failed to run servers:", err)
	}
}

func runServers(ctx context.Context, cfg *Config, userServer userv1.UserServiceServer, authMiddleware *auth.AuthMiddleware) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// gRPC Server
	g.Go(func() error {
		return runGRPCServer(ctx, cfg.GRPCPort, userServer, authMiddleware)
	})

	// HTTP Gateway
	g.Go(func() error {
		return runHTTPServer(ctx, cfg.HTTPPort, cfg.GRPCPort, authMiddleware)
	})

	// Graceful shutdown
	g.Go(func() error {
		return handleShutdown(ctx, cancel)
	})

	return g.Wait()
}

func runGRPCServer(ctx context.Context, port string, userServer userv1.UserServiceServer, authMiddleware *auth.AuthMiddleware) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authMiddleware.GRPCAuthInterceptor,
		),
	)

	userv1.RegisterUserServiceServer(server, userServer)

	log.Printf("gRPC server listening on :%s", port)
	return server.Serve(lis)
}

func handleShutdown(ctx context.Context, cancel context.CancelFunc) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case sig := <-sigCh:
		log.Printf("Received signal: %s", sig)
		cancel()
		// Give some time for graceful shutdown
		time.Sleep(5 * time.Second)
		return nil
	}
}
