package app

import (
	"context"
	"log"
	"net"

	genprotov1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"google.golang.org/grpc"
)

func (a *App) runGRPC(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":"+a.cfg.GRPCPort)
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			a.middleware.Auth.GRPCAuthInterceptor,
		),
	)

	genprotov1.RegisterUserServiceServer(server, a.serviceServer.userServer)
	genprotov1.RegisterOAuthServiceServer(server, a.serviceServer.oauthServer)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
		lis.Close()
	}()

	log.Printf("gRPC server running on :%s", a.cfg.GRPCPort)
	return server.Serve(lis)
}
