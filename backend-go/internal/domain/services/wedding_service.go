package services

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
)

type WeddingService struct {
	weddings  repositories.WeddingRepository
	guests    repositories.GuestRepository
	templates repositories.TemplateRepository
	payments  repositories.PaymentRepository
	subs      repositories.SubscriptionRepository
}

func NewWeddingService(w repositories.WeddingRepository, g repositories.GuestRepository, t repositories.TemplateRepository, p repositories.PaymentRepository, s repositories.SubscriptionRepository) *WeddingService {
	return &WeddingService{weddings: w, guests: g, templates: t, payments: p, subs: s}
}

func (s *WeddingService) Create(ctx context.Context, userID uuid.UUID, templateID *uuid.UUID, config map[string]any) (*entities.Wedding, error) {
	w := &entities.Wedding{UserID: userID, TemplateID: templateID, Status: entities.WeddingDraft, ConfigData: config}
	return s.weddings.Create(ctx, w)
}

func (s *WeddingService) GetByID(ctx context.Context, weddingID uuid.UUID) (*entities.Wedding, error) {
	return s.weddings.GetByID(ctx, weddingID)
}

func (s *WeddingService) SetTemplate(ctx context.Context, weddingID uuid.UUID, templateID uuid.UUID) (*entities.Wedding, error) {
	return s.weddings.SetTemplate(ctx, weddingID, templateID)
}

func (s *WeddingService) UpdateConfig(ctx context.Context, weddingID uuid.UUID, config map[string]any) (*entities.Wedding, error) {
	return s.weddings.UpdateConfig(ctx, weddingID, config)
}

func (s *WeddingService) SetDomain(ctx context.Context, weddingID uuid.UUID, domain string) (*entities.Wedding, error) {
	d := strings.TrimSpace(domain)
	if d == "" {
		return nil, errors.New("domain required")
	}
	return s.weddings.SetDomain(ctx, weddingID, d)
}

func (s *WeddingService) SetSlug(ctx context.Context, weddingID uuid.UUID, slug string) (*entities.Wedding, error) {
	sl := strings.TrimSpace(slug)
	if sl == "" {
		return nil, errors.New("slug required")
	}
	return s.weddings.SetSlug(ctx, weddingID, sl)
}

func (s *WeddingService) SetPayment(ctx context.Context, weddingID uuid.UUID, paymentID uuid.UUID) (*entities.Wedding, error) {
	return s.weddings.SetPayment(ctx, weddingID, paymentID)
}

func (s *WeddingService) Publish(ctx context.Context, weddingID uuid.UUID, userID uuid.UUID) (*entities.Wedding, error) {
	sub, _ := s.subs.GetByUser(ctx, userID)
	if sub != nil && strings.ToLower(sub.Status) == "active" {
		return s.weddings.Publish(ctx, weddingID)
	}
	w, err := s.weddings.GetByID(ctx, weddingID)
	if err != nil {
		return nil, err
	}
	if w.PaymentID == nil {
		return nil, errors.New("payment required")
	}
	pay, err := s.payments.GetByID(ctx, *w.PaymentID)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(string(pay.Status)) != "paid" {
		return nil, errors.New("payment not paid")
	}
	return s.weddings.Publish(ctx, weddingID)
}

func (s *WeddingService) AddGuest(ctx context.Context, g *entities.Guest) (*entities.Guest, error) {
	return s.guests.Add(ctx, g)
}

func (s *WeddingService) UpdateGuest(ctx context.Context, g *entities.Guest) (*entities.Guest, error) {
	return s.guests.Update(ctx, g)
}

func (s *WeddingService) DeleteGuest(ctx context.Context, id uuid.UUID) error {
	return s.guests.Delete(ctx, id)
}

func (s *WeddingService) ListGuests(ctx context.Context, weddingID uuid.UUID) ([]*entities.Guest, error) {
	return s.guests.ListByWedding(ctx, weddingID)
}
