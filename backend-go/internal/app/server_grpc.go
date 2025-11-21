package app

import (
	"context"
	"log"
	"net"

	genproto "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto"
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

	genproto.RegisterUserServiceServer(server, a.serviceServer.userServer)

	log.Printf("gRPC server running on :%s", a.cfg.GRPCPort)
	return server.Serve(lis)
}
