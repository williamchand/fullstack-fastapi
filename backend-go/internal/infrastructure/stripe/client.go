package stripeinfra

import (
    "github.com/stripe/stripe-go/v78"
    "github.com/stripe/stripe-go/v78/checkout/session"
    "github.com/stripe/stripe-go/v78/price"
    "github.com/stripe/stripe-go/v78/webhook"
)

type Client struct{}

func New(secret string) *Client {
	stripe.Key = secret
	return &Client{}
}

func (c *Client) CreateCheckoutSession(priceID, successURL, cancelURL string, metadata map[string]string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			Price:    stripe.String(priceID),
			Quantity: stripe.Int64(1),
		}},
	}
	if metadata != nil {
		params.Metadata = metadata
	}
	s, err := session.New(params)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (c *Client) ConstructEvent(payload []byte, sig string, webhookSecret string) (stripe.Event, error) {
    return webhook.ConstructEvent(payload, sig, webhookSecret)
}

func (c *Client) GetPrice(priceID string) (*stripe.Price, error) {
	p, err := price.Get(priceID, nil)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (c *Client) GetCheckoutSession(id string) (*stripe.CheckoutSession, error) {
    s, err := session.Get(id, nil)
    if err != nil {
        return nil, err
    }
    return s, nil
}
