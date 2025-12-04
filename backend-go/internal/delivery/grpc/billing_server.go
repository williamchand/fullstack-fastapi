package grpc

import (
    "context"

    salonappv1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
    "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/auth"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
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
	if err != nil {
		return nil, err
	}
	return &salonappv1.CreateCheckoutSessionResponse{Url: url, SessionId: id}, nil
}

func (s *billingServer) GetSubscriptionStatus(ctx context.Context, _ *emptypb.Empty) (*salonappv1.GetSubscriptionStatusResponse, error) {
	user := auth.UserFromContext(ctx)
	sub, err := s.svc.GetSubscriptionStatus(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return &salonappv1.GetSubscriptionStatusResponse{Status: sub.Status, StripeSubscriptionId: str(sub.StripeSubscriptionID)}, nil
}

func (s *billingServer) HandleWebhook(ctx context.Context, req *salonappv1.HandleWebhookRequest) (*emptypb.Empty, error) {
    switch req.Provider {
    case salonappv1.Provider_PROVIDER_STRIPE:
        if err := s.svc.HandleWebhook(ctx, req.Payload, req.Signature); err != nil { return nil, err }
        return &emptypb.Empty{}, nil
    case salonappv1.Provider_PROVIDER_DOKU:
        // Parse payload JSON for DOKU notification
        var invoice, session, statusStr, currency, amount string
        // Lightweight decode: try known fields
        type dokuPayload struct {
            Response struct {
                Order struct { InvoiceNumber string `json:"invoice_number"`; SessionID string `json:"session_id"`; Currency string `json:"currency"` } `json:"order"`
            } `json:"response"`
            Order struct { InvoiceNumber string `json:"invoice_number"` } `json:"order"`
            Status string `json:"status"`
            Currency string `json:"currency"`
            Amount string `json:"amount"`
            SessionID string `json:"session_id"`
        }
        var dp dokuPayload
        if err := json.Unmarshal(req.Payload, &dp); err == nil {
            invoice = dp.Response.Order.InvoiceNumber
            if invoice == "" { invoice = dp.Order.InvoiceNumber }
            session = dp.Response.Order.SessionID
            if session == "" { session = dp.SessionID }
            currency = dp.Response.Order.Currency
            if currency == "" { currency = dp.Currency }
            amount = dp.Amount
            statusStr = dp.Status
        }
        if err := s.svc.HandleDokuNotification(ctx, invoice, session, statusStr, currency, amount); err != nil { return nil, err }
        return &emptypb.Empty{}, nil
    default:
        return nil, status.Error(codes.InvalidArgument, "unsupported provider")
    }
}

// CreateDokuPayment initiates a Jokul Checkout payment and returns payment URL
func (s *billingServer) CreateDokuPayment(ctx context.Context, req *salonappv1.CreateDokuPaymentRequest) (*salonappv1.CreateDokuPaymentResponse, error) {
    user := auth.UserFromContext(ctx)
    amount := req.AmountIdr
    invoice := req.InvoiceNumber
    due := int(req.PaymentDueMinutes)
    if amount <= 0 || invoice == "" {
        return nil, status.Error(codes.InvalidArgument, "invalid amount or invoice number")
    }
    url, txid, err := s.svc.CreateDokuPayment(ctx, user.ID, amount, invoice, due)
    if err != nil {
        return nil, err
    }
    return &salonappv1.CreateDokuPaymentResponse{PaymentUrl: url, TransactionId: txid}, nil
}

// statusError wraps a simple string into an error; reuse existing error handling
func (s *billingServer) RefreshPaymentStatus(ctx context.Context, req *salonappv1.RefreshPaymentStatusRequest) (*salonappv1.PaymentStatusResponse, error) {
    user := auth.UserFromContext(ctx)
    if req.TransactionId == "" || req.Provider == salonappv1.Provider_PROVIDER_UNSPECIFIED {
        return nil, status.Error(codes.InvalidArgument, "transaction_id and provider are required")
    }
    var providerStr string
    switch req.Provider {
    case salonappv1.Provider_PROVIDER_STRIPE:
        providerStr = "stripe"
    case salonappv1.Provider_PROVIDER_DOKU:
        providerStr = "doku"
    default:
        return nil, status.Error(codes.InvalidArgument, "unsupported provider")
    }
    statusStr, err := s.svc.RefreshPaymentStatus(ctx, user.ID, req.TransactionId, providerStr)
    if err != nil {
        return nil, err
    }
    return &salonappv1.PaymentStatusResponse{TransactionId: req.TransactionId, Provider: req.Provider, Status: statusStr}, nil
}

func (s *billingServer) CheckDailySubscriptions(ctx context.Context, _ *emptypb.Empty) (*salonappv1.CheckDailySubscriptionsResponse, error) {
    updated, expired, err := s.svc.CheckDailySubscriptions(ctx)
    if err != nil {
        return nil, err
    }
    return &salonappv1.CheckDailySubscriptionsResponse{Updated: int32(updated), Expired: int32(expired)}, nil
}

func str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
