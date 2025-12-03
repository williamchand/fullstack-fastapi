package services

import (
    "context"
    "encoding/json"
    "fmt"

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
	case "checkout.session.completed":
		var s stripe.CheckoutSession
		err := json.Unmarshal(evt.Data.Raw, &s)
		if err != nil {
			return err
		}
		// Update payment status
		_, _ = b.payRepo.UpdateStatus(ctx, s.ID, entities.PaymentStatusPaid, nil, nil, map[string]any{"customer": s.Customer.ID})
	case "checkout.session.expired":
		var s stripe.CheckoutSession
		err := json.Unmarshal(evt.Data.Raw, &s)
		if err != nil {
			return err
		}
		_, _ = b.payRepo.UpdateStatus(ctx, s.ID, entities.PaymentStatusFailed, nil, nil, map[string]any{"reason": "expired"})
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
