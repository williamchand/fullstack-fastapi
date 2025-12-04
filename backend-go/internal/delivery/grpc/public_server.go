package grpc

import (
	"context"

	salonappv1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
)

type publicServer struct {
	salonappv1.UnimplementedPublicServiceServer
	svc *services.PublicService
}

func NewPublicServer(svc *services.PublicService) salonappv1.PublicServiceServer {
	return &publicServer{svc: svc}
}

func (s *publicServer) GetRegions(ctx context.Context, req *salonappv1.GetRegionsRequest) (*salonappv1.GetRegionsResponse, error) {
	details := s.svc.GetRegions(req.Regions)
	resp := &salonappv1.GetRegionsResponse{Regions: make([]*salonappv1.RegionInfo, 0, len(details))}
	for _, d := range details {
		resp.Regions = append(resp.Regions, &salonappv1.RegionInfo{Region: d.Region, CountryCode: d.CountryCode, Supported: d.Supported})
	}
	return resp, nil
}
