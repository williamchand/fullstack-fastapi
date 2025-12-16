package grpc

import (
	"context"

	"github.com/google/uuid"
	salonappv1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/util"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type dataSourceServer struct {
	salonappv1.UnimplementedDataSourceServiceServer
	svc *services.DataSourceService
}

func NewDataSourceServer(svc *services.DataSourceService) salonappv1.DataSourceServiceServer {
	return &dataSourceServer{svc: svc}
}

func (s *dataSourceServer) CreateDataSource(ctx context.Context, req *salonappv1.CreateDataSourceRequest) (*salonappv1.CreateDataSourceResponse, error) {
	user := util.UserFromContext(ctx)
	ds, err := s.svc.CreateDataSource(ctx, user.ID, req.Name, req.Type, req.Host, req.Port, req.DatabaseName, req.Username, req.Password, map[string]any{})
	if err != nil {
		return nil, err
	}
	return &salonappv1.CreateDataSourceResponse{DataSource: s.toProto(ds)}, nil
}

func (s *dataSourceServer) TestConnection(ctx context.Context, req *salonappv1.TestConnectionRequest) (*salonappv1.TestConnectionResponse, error) {
	user := util.UserFromContext(ctx)
	ds, err := s.svc.GetByID(ctx, uuid.MustParse(req.Id), user.ID)
	if err != nil {
		return nil, err
	}
	ok, msg := s.svc.TestConnection(ctx, ds)
	return &salonappv1.TestConnectionResponse{Ok: ok, Message: msg}, nil
}

func (s *dataSourceServer) IntrospectSchema(ctx context.Context, req *salonappv1.IntrospectSchemaRequest) (*salonappv1.IntrospectSchemaResponse, error) {
	user := util.UserFromContext(ctx)
	ds, err := s.svc.GetByID(ctx, uuid.MustParse(req.Id), user.ID)
	if err != nil {
		return nil, err
	}
	schema, summary, err := s.svc.IntrospectSchema(ctx, ds)
	if err != nil {
		return nil, err
	}
	return &salonappv1.IntrospectSchemaResponse{SchemaJson: schema, Summary: summary}, nil
}

func (s *dataSourceServer) SetAICredential(ctx context.Context, req *salonappv1.SetAICredentialRequest) (*emptypb.Empty, error) {
	user := util.UserFromContext(ctx)
	if err := s.svc.SetAICredential(ctx, user.ID, req.Provider, req.ApiKey); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *dataSourceServer) AskQuestion(ctx context.Context, req *salonappv1.AskQuestionRequest) (*salonappv1.AskQuestionResponse, error) {
	user := util.UserFromContext(ctx)
	ds, err := s.svc.GetByID(ctx, uuid.MustParse(req.DataSourceId), user.ID)
	if err != nil {
		return nil, err
	}
	ans, sql, res, err := s.svc.AskQuestion(ctx, ds, "openai", req.Question)
	if err != nil {
		return nil, err
	}
	return &salonappv1.AskQuestionResponse{Answer: ans, Sql: sql, ResultJson: res}, nil
}

func (s *dataSourceServer) toProto(ds *entities.DataSource) *salonappv1.DataSource {
	return &salonappv1.DataSource{
		Id:           ds.ID.String(),
		Name:         ds.Name,
		Type:         ds.Type,
		Host:         ds.Host,
		Port:         ds.Port,
		DatabaseName: ds.DatabaseName,
		Username:     ds.Username,
		CreatedAt:    timestamppb.New(ds.CreatedAt),
		UpdatedAt:    timestamppb.New(ds.UpdatedAt),
	}
}
