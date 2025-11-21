package app

import (
	userv1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/delivery/grpc"
)

type ServiceServer struct {
	userServer userv1.UserServiceServer
}

func initServiceServer(appServices *AppServices) *ServiceServer {
	userServer := grpc.NewUserServer(appServices.UserService)
	return &ServiceServer{
		userServer: userServer,
	}
}
