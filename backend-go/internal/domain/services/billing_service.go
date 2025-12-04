package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"

	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/repositories"
	stripeinfra "github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/stripe"
)

type BillingService struct {
	cfg     *config.Config
	subs    repositories.SubscriptionRepository
	payRepo repositories.PaymentRepository
	stripe  *stripeinfra.Client
	doku    repositories.DokuClient
}

func NewBillingService(cfg *config.Config, subs repositories.SubscriptionRepository, pay repositories.PaymentRepository, stripeClient *stripeinfra.Client, dokuClient repositories.DokuClient) *BillingService {
	return &BillingService{cfg: cfg, subs: subs, payRepo: pay, stripe: stripeClient, doku: dokuClient}
}

func (b *BillingService) CreateCheckoutSession(ctx context.Context, userID uuid.UUID, successURL, cancelURL string) (string, string, error) {
	s, err := b.stripe.CreateCheckoutSession(b.cfg.Stripe.PriceID, successURL, cancelURL, map[string]string{"user_id": userID.String()})
	if err != nil {
		return "", "", err
	}
	// Get price info for amount/currency
	pr, err := b.stripe.GetPrice(b.cfg.Stripe.PriceID)
	if err != nil {
		return "", "", err
	}
	amount := float64(pr.UnitAmount) / 100.0
	currency := string(pr.Currency)
	_, err = b.payRepo.Create(ctx, &entities.Payment{
		UserID:        userID,
		Provider:      entities.PaymentProviderStripe,
		Amount:        amount,
		Currency:      currency,
		Status:        entities.PaymentStatusPending,
		TransactionID: s.ID,
		ExtraMetadata: map[string]any{"checkout_url": s.URL},
	})
	if err != nil {
		return "", "", err
	}
	return s.URL, s.ID, nil
}

func (b *BillingService) HandleWebhook(ctx context.Context, payload []byte, sig string) error {
	evt, err := b.stripe.ConstructEvent(payload, sig, b.cfg.Stripe.WebhookSecret)
	if err != nil {
		return err
	}
	switch evt.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var s stripe.CheckoutSession
		err := json.Unmarshal(evt.Data.Raw, &s)
		if err != nil {
			return err
		}
		// Update payment status
		_, _ = b.payRepo.UpdateStatus(ctx, s.ID, entities.PaymentStatusPaid, nil, nil, map[string]any{"customer": s.Customer.ID})
	case stripe.EventTypeCheckoutSessionExpired:
		var s stripe.CheckoutSession
		err := json.Unmarshal(evt.Data.Raw, &s)
		if err != nil {
			return err
		}
		_, _ = b.payRepo.UpdateStatus(ctx, s.ID, entities.PaymentStatusFailed, nil, nil, map[string]any{"reason": entities.PaymentStatusExpired})
	}
	return nil
}

func (b *BillingService) GetSubscriptionStatus(ctx context.Context, userID uuid.UUID) (*entities.Subscription, error) {
	return b.subs.GetByUser(ctx, userID)
}

// CreateDokuPayment initiates a DOKU Jokul checkout and records a pending payment
func (b *BillingService) CreateDokuPayment(ctx context.Context, userID uuid.UUID, amountIDR int64, invoiceNumber string, paymentDueMinutes int) (string, string, error) {
	if b.doku == nil {
		return "", "", fmt.Errorf("doku client not configured")
	}
	url, txid, amt, currency, err := b.doku.CreatePayment(ctx, amountIDR, invoiceNumber, paymentDueMinutes)
	if err != nil {
		return "", "", err
	}
	amount := float64(amt) // IDR has no decimals
	_, err = b.payRepo.Create(ctx, &entities.Payment{
		UserID:        userID,
		Provider:      entities.PaymentProviderDoku,
		Amount:        amount,
		Currency:      currency,
		Status:        entities.PaymentStatusPending,
		TransactionID: txid,
		ExtraMetadata: map[string]any{"payment_url": url, "invoice_number": invoiceNumber},
	})
	if err != nil {
		return "", "", err
	}
	return url, txid, nil
}

// HandleDokuNotification processes DOKU HTTP notification and updates payment status
func (b *BillingService) HandleDokuNotification(ctx context.Context, invoiceNumber, sessionID, currency, amountStr string, statusStr entities.PaymentStatus) error {
	txid := sessionID
	if txid == "" {
		txid = invoiceNumber
	}
	if txid == "" {
		return fmt.Errorf("missing transaction identifier")
	}
	var amount *float64
	if amountStr != "" {
		// DOKU sends amount as string, convert to float
		var a float64
		_, err := fmt.Sscanf(amountStr, "%f", &a)
		if err == nil {
			amount = &a
		}
	}
	var curr *string
	if currency != "" {
		curr = &currency
	}
	// Map status
	st := entities.PaymentStatusPending
	switch statusStr {
	case entities.PaymentStatusSuccess, entities.PaymentStatusPaid, entities.PaymentStatusCompleted:
		st = entities.PaymentStatusPaid
	case entities.PaymentStatusFailed, entities.PaymentStatusExpired:
		st = entities.PaymentStatusFailed
	}
	_, err := b.payRepo.UpdateStatus(ctx, txid, st, amount, curr, map[string]any{"provider": entities.PaymentProviderDoku})
	return err
}

// RefreshPaymentStatus checks provider for latest status and updates payment; returns status string
func (b *BillingService) RefreshPaymentStatus(ctx context.Context, userID uuid.UUID, txid string, provider entities.PaymentProvider) (string, error) {
	p, err := b.payRepo.GetByTransaction(ctx, txid)
	if err != nil {
		return "", err
	}
	if p.UserID != userID {
		return "", fmt.Errorf("not owner")
	}
	switch provider {
	case entities.PaymentProviderStripe:
		sess, err := b.stripe.GetCheckoutSession(txid)
		if err != nil {
			return "", err
		}
		var newStatus entities.PaymentStatus
		switch sess.Status {
		case stripe.CheckoutSessionStatusComplete:
			newStatus = entities.PaymentStatusPaid
			_, _ = b.payRepo.UpdateStatus(ctx, txid, newStatus, nil, nil, map[string]any{"provider": entities.PaymentProviderStripe})
		case stripe.CheckoutSessionStatusExpired:
			newStatus = entities.PaymentStatusFailed
			_, _ = b.payRepo.UpdateStatus(ctx, txid, newStatus, nil, nil, map[string]any{"provider": entities.PaymentProviderStripe, "reason": entities.PaymentStatusExpired})
		default:
			// pending
			newStatus = entities.PaymentStatusPending
		}
		return string(newStatus), nil
	case entities.PaymentProviderDoku:
		// DOKU doesn't provide a simple status fetch in Jokul; rely on notification
		// Here we just return current stored status
		return string(p.Status), nil
	}
	return string(p.Status), nil
}

// CheckDailySubscriptions scans and updates subscription statuses based on period end
func (b *BillingService) CheckDailySubscriptions(ctx context.Context) (updated int, expired int, err error) {
	subs, err := b.subs.ListAll(ctx)
	if err != nil {
		return 0, 0, err
	}
	now := time.Now()
	for _, s := range subs {
		if s.CurrentPeriodEnd != nil && now.After(*s.CurrentPeriodEnd) && s.Status != entities.PaymentStatusExpired {
			s.Status = entities.PaymentStatusExpired
			_, _ = b.subs.Upsert(ctx, s)
			expired++
		} else {
			updated++
		}
	}
	return
}
