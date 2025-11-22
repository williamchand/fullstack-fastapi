package app

import (
	genprotov1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/delivery/grpc"
)

type ServiceServer struct {
	userServer  genprotov1.UserServiceServer
	oauthServer genprotov1.OAuthServiceServer
}

func initServiceServer(appServices *AppServices) *ServiceServer {
	userServer := grpc.NewUserServer(appServices.UserService)
	oauthServer := grpc.NewOAuthServer(appServices.OauthService)
	return &ServiceServer{
		userServer:  userServer,
		oauthServer: oauthServer,
	}
}
