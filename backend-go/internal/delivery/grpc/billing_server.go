package grpc

import (
    "context"
    salonappv1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/auth"
    "google.golang.org/protobuf/types/known/emptypb"
)

type billingServer struct {
    salonappv1.UnimplementedBillingServiceServer
    svc *services.BillingService
}

func NewBillingServer(svc *services.BillingService) salonappv1.BillingServiceServer {
    return &billingServer{svc: svc}
}

func (s *billingServer) CreateCheckoutSession(ctx context.Context, req *salonappv1.CreateCheckoutSessionRequest) (*salonappv1.CreateCheckoutSessionResponse, error) {
    user := auth.UserFromContext(ctx)
    url, id, err := s.svc.CreateCheckoutSession(ctx, user.ID, req.SuccessUrl, req.CancelUrl)
    if err != nil { return nil, err }
    return &salonappv1.CreateCheckoutSessionResponse{Url: url, SessionId: id}, nil
}

func (s *billingServer) GetSubscriptionStatus(ctx context.Context, _ *emptypb.Empty) (*salonappv1.GetSubscriptionStatusResponse, error) {
    user := auth.UserFromContext(ctx)
    sub, err := s.svc.GetSubscriptionStatus(ctx, user.ID)
    if err != nil { return nil, err }
    var ts *int64
    if sub.CurrentPeriodEnd != nil { t := sub.CurrentPeriodEnd.Unix(); ts = &t }
    return &salonappv1.GetSubscriptionStatusResponse{Status: sub.Status, StripeSubscriptionId: str(sub.StripeSubscriptionID)}, nil
}

func (s *billingServer) HandleWebhook(ctx context.Context, req *salonappv1.HandleWebhookRequest) (*emptypb.Empty, error) {
    if err := s.svc.HandleWebhook(ctx, req.Payload, req.Signature); err != nil { return nil, err }
    return &emptypb.Empty{}, nil
}

func str(p *string) string { if p==nil {return ""}; return *p }
