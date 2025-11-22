package app

import (
	genprotov1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/delivery/grpc"
)

type ServiceServer struct {
	userServer genprotov1.UserServiceServer
}

func initServiceServer(appServices *AppServices) *ServiceServer {
	userServer := grpc.NewUserServer(appServices.UserService)
	return &ServiceServer{
		userServer: userServer,
	}
}
