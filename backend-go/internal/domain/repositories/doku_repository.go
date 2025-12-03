package repositories

import "context"

// DokuClient defines Jokul Checkout API operations needed by our app
type DokuClient interface {
    // CreatePayment initiates a DOKU payment and returns the payment URL and transaction ID
    // amountIDR: integer amount in IDR (no decimals)
    // invoiceNumber: merchant invoice identifier
    // paymentDueMinutes: expiration time in minutes
    CreatePayment(ctx context.Context, amountIDR int64, invoiceNumber string, paymentDueMinutes int) (paymentURL string, transactionID string, idrAmount int64, currency string, err error)
}

