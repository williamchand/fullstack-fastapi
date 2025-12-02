package grpc

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	salonappv1 "github.com/williamchand/fullstack-fastapi/backend-go/gen/proto/v1"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/services"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

type weddingServer struct {
	salonappv1.UnimplementedWeddingServiceServer
	svc *services.WeddingService
}

func NewWeddingServer(svc *services.WeddingService) salonappv1.WeddingServiceServer {
	return &weddingServer{svc: svc}
}

func (s *weddingServer) CreateWedding(ctx context.Context, req *salonappv1.CreateWeddingRequest) (*salonappv1.WeddingResponse, error) {
	user := auth.UserFromContext(ctx)
	var cfg map[string]any
	if req.ConfigJson != "" {
		_ = json.Unmarshal([]byte(req.ConfigJson), &cfg)
	}
	var tmplID *uuid.UUID
	if req.TemplateId != "" {
		id := uuid.MustParse(req.TemplateId)
		tmplID = &id
	}
	w, err := s.svc.Create(ctx, user.ID, tmplID, cfg)
	if err != nil {
		return nil, err
	}
	return &salonappv1.WeddingResponse{Wedding: toWeddingProto(w)}, nil
}

func (s *weddingServer) GetWedding(ctx context.Context, req *salonappv1.GetWeddingRequest) (*salonappv1.WeddingResponse, error) {
	w, err := s.svc.GetByID(ctx, uuid.MustParse(req.Id))
	if err != nil {
		return nil, err
	}
	return &salonappv1.WeddingResponse{Wedding: toWeddingProto(w)}, nil
}

func (s *weddingServer) UpdateConfig(ctx context.Context, req *salonappv1.UpdateConfigRequest) (*salonappv1.WeddingResponse, error) {
	var cfg map[string]any
	if req.ConfigJson != "" {
		_ = json.Unmarshal([]byte(req.ConfigJson), &cfg)
	}
	w, err := s.svc.UpdateConfig(ctx, uuid.MustParse(req.Id), cfg)
	if err != nil {
		return nil, err
	}
	return &salonappv1.WeddingResponse{Wedding: toWeddingProto(w)}, nil
}

func (s *weddingServer) SetTemplate(ctx context.Context, req *salonappv1.SetTemplateRequest) (*salonappv1.WeddingResponse, error) {
	w, err := s.svc.SetTemplate(ctx, uuid.MustParse(req.Id), uuid.MustParse(req.TemplateId))
	if err != nil {
		return nil, err
	}
	return &salonappv1.WeddingResponse{Wedding: toWeddingProto(w)}, nil
}

func (s *weddingServer) SetDomain(ctx context.Context, req *salonappv1.SetDomainRequest) (*salonappv1.WeddingResponse, error) {
	w, err := s.svc.SetDomain(ctx, uuid.MustParse(req.Id), req.Domain)
	if err != nil {
		return nil, err
	}
	return &salonappv1.WeddingResponse{Wedding: toWeddingProto(w)}, nil
}

func (s *weddingServer) SetSlug(ctx context.Context, req *salonappv1.SetSlugRequest) (*salonappv1.WeddingResponse, error) {
	w, err := s.svc.SetSlug(ctx, uuid.MustParse(req.Id), req.Slug)
	if err != nil {
		return nil, err
	}
	return &salonappv1.WeddingResponse{Wedding: toWeddingProto(w)}, nil
}

func (s *weddingServer) Publish(ctx context.Context, req *salonappv1.PublishRequest) (*salonappv1.WeddingResponse, error) {
	user := auth.UserFromContext(ctx)
	w, err := s.svc.Publish(ctx, uuid.MustParse(req.Id), user.ID)
	if err != nil {
		return nil, err
	}
	return &salonappv1.WeddingResponse{Wedding: toWeddingProto(w)}, nil
}

func (s *weddingServer) AddGuest(ctx context.Context, req *salonappv1.AddGuestRequest) (*salonappv1.GuestResponse, error) {
	g, err := s.svc.AddGuest(ctx, &entities.Guest{WeddingID: uuid.MustParse(req.Id), Name: req.Name, Contact: req.Contact, RSVPStatus: entities.RSVPStatus(req.RsvpStatus), Message: strPtr(req.Message)})
	if err != nil {
		return nil, err
	}
	return &salonappv1.GuestResponse{Guest: toGuestProto(g)}, nil
}

func (s *weddingServer) UpdateGuest(ctx context.Context, req *salonappv1.UpdateGuestRequest) (*salonappv1.GuestResponse, error) {
	g, err := s.svc.UpdateGuest(ctx, &entities.Guest{ID: uuid.MustParse(req.GuestId), WeddingID: uuid.MustParse(req.Id), Name: req.Name, Contact: req.Contact, RSVPStatus: entities.RSVPStatus(req.RsvpStatus), Message: strPtr(req.Message)})
	if err != nil {
		return nil, err
	}
	return &salonappv1.GuestResponse{Guest: toGuestProto(g)}, nil
}

func (s *weddingServer) DeleteGuest(ctx context.Context, req *salonappv1.DeleteGuestRequest) (*emptypb.Empty, error) {
	err := s.svc.DeleteGuest(ctx, uuid.MustParse(req.GuestId))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *weddingServer) ListGuests(ctx context.Context, req *salonappv1.ListGuestsRequest) (*salonappv1.GuestsResponse, error) {
	rows, err := s.svc.ListGuests(ctx, uuid.MustParse(req.Id))
	if err != nil {
		return nil, err
	}
	res := &salonappv1.GuestsResponse{Guests: make([]*salonappv1.Guest, 0, len(rows))}
	for _, g := range rows {
		res.Guests = append(res.Guests, toGuestProto(g))
	}
	return res, nil
}

func toWeddingProto(w *entities.Wedding) *salonappv1.Wedding {
	var tmplID, payID, domain, slug string
	if w.TemplateID != nil {
		tmplID = w.TemplateID.String()
	}
	if w.PaymentID != nil {
		payID = w.PaymentID.String()
	}
	if w.CustomDomain != nil {
		domain = *w.CustomDomain
	}
	if w.Slug != nil {
		slug = *w.Slug
	}
	b, _ := json.Marshal(w.ConfigData)
	return &salonappv1.Wedding{Id: w.ID.String(), UserId: w.UserID.String(), TemplateId: tmplID, PaymentId: payID, Status: string(w.Status), CustomDomain: domain, Slug: slug, ConfigJson: string(b), CreatedAt: w.CreatedAt.Unix()}
}

func toGuestProto(g *entities.Guest) *salonappv1.Guest {
	var msg string
	if g.Message != nil {
		msg = *g.Message
	}
	return &salonappv1.Guest{Id: g.ID.String(), WeddingId: g.WeddingID.String(), Name: g.Name, Contact: g.Contact, RsvpStatus: string(g.RSVPStatus), Message: msg, CreatedAt: g.CreatedAt.Unix()}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
