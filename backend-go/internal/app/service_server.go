package app

import (
	genprotov1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/delivery/grpc"
)

type ServiceServer struct {
	userServer    genprotov1.UserServiceServer
	oauthServer   genprotov1.OAuthServiceServer
	billingServer genprotov1.BillingServiceServer
}

func initServiceServer(appServices *AppServices) *ServiceServer {
	userServer := grpc.NewUserServer(appServices.UserService)
	oauthServer := grpc.NewOAuthServer(appServices.OauthService)
	billServer := grpc.NewBillingServer(appServices.BillingService)
	return &ServiceServer{
		userServer:    userServer,
		oauthServer:   oauthServer,
		billingServer: billServer,
	}
}
