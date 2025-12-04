package app

import (
	genprotov1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/delivery/grpc"
)

type ServiceServer struct {
	userServer       genprotov1.UserServiceServer
	oauthServer      genprotov1.OAuthServiceServer
	dataSourceServer genprotov1.DataSourceServiceServer
	billingServer    genprotov1.BillingServiceServer
	weddingServer    genprotov1.WeddingServiceServer
	publicServer     genprotov1.PublicServiceServer
}

func initServiceServer(appServices *AppServices) *ServiceServer {
	userServer := grpc.NewUserServer(appServices.UserService)
	oauthServer := grpc.NewOAuthServer(appServices.OauthService)
	dsServer := grpc.NewDataSourceServer(appServices.DataSourceService)
	billServer := grpc.NewBillingServer(appServices.BillingService)
	wedServer := grpc.NewWeddingServer(appServices.WeddingService)
	pubServer := grpc.NewPublicServer(appServices.PublicService)
	return &ServiceServer{
		userServer:       userServer,
		oauthServer:      oauthServer,
		dataSourceServer: dsServer,
		billingServer:    billServer,
		weddingServer:    wedServer,
		publicServer:     pubServer,
	}
}
