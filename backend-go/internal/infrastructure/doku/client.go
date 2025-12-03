package doku

import (
    "bytes"
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/google/uuid"
)

type Client struct {
    baseURL    string
    clientID   string
    secretKey  string
    httpClient *http.Client
}

func New(baseURL, clientID, secretKey string) *Client {
    return &Client{
        baseURL:    strings.TrimRight(baseURL, "/"),
        clientID:   clientID,
        secretKey:  secretKey,
        httpClient: &http.Client{Timeout: 15 * time.Second},
    }
}

type basicRequest struct {
    Order struct {
        Amount        int64  `json:"amount"`
        InvoiceNumber string `json:"invoice_number"`
    } `json:"order"`
    Payment struct {
        PaymentDueDate int `json:"payment_due_date"`
    } `json:"payment"`
}

type dokuResponse struct {
    Message  []string `json:"message"`
    Response struct {
        Order struct {
            Amount        string `json:"amount"`
            InvoiceNumber string `json:"invoice_number"`
            Currency      string `json:"currency"`
            SessionID     string `json:"session_id"`
        } `json:"order"`
        Payment struct {
            URL        string `json:"url"`
            PaymentURL string `json:"payment_url"`
        } `json:"payment"`
    } `json:"response"`
}

// CreatePayment performs Jokul Checkout initiation and returns URL and transaction id
func (c *Client) CreatePayment(ctx context.Context, amountIDR int64, invoiceNumber string, paymentDueMinutes int) (string, string, int64, string, error) {
    // Build body
    reqBody := basicRequest{}
    reqBody.Order.Amount = amountIDR
    reqBody.Order.InvoiceNumber = invoiceNumber
    reqBody.Payment.PaymentDueDate = paymentDueMinutes
    bodyBytes, err := json.Marshal(reqBody)
    if err != nil {
        return "", "", 0, "", fmt.Errorf("doku: marshal body failed: %w", err)
    }

    // Headers
    requestID := uuid.NewString()
    timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
    path := "/checkout/v1/payment"

    // Digest
    digestHash := sha256.Sum256(bodyBytes)
    digest := "SHA-256=" + base64.StdEncoding.EncodeToString(digestHash[:])

    // Signature base string (per Jokul docs)
    baseString := strings.Join([]string{
        "Client-Id:" + c.clientID,
        "Request-Id:" + requestID,
        "Request-Timestamp:" + timestamp,
        "Request-Target:" + path,
        "Digest:" + digest,
    }, "\n")
    mac := hmac.New(sha256.New, []byte(c.secretKey))
    mac.Write([]byte(baseString))
    signature := "HMACSHA256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))

    // Request
    url := c.baseURL + path
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
    if err != nil {
        return "", "", 0, "", fmt.Errorf("doku: request build failed: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Client-Id", c.clientID)
    req.Header.Set("Request-Id", requestID)
    req.Header.Set("Request-Timestamp", timestamp)
    req.Header.Set("Signature", signature)

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return "", "", 0, "", fmt.Errorf("doku: request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 {
        return "", "", 0, "", fmt.Errorf("doku: non-2xx status %d", resp.StatusCode)
    }

    var dr dokuResponse
    if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
        return "", "", 0, "", fmt.Errorf("doku: decode response failed: %w", err)
    }

    // Extract payment URL and transaction id
    paymentURL := dr.Response.Payment.URL
    if paymentURL == "" {
        paymentURL = dr.Response.Payment.PaymentURL
    }
    txid := dr.Response.Order.SessionID
    if txid == "" {
        txid = dr.Response.Order.InvoiceNumber
    }
    // Default currency
    currency := dr.Response.Order.Currency
    if currency == "" {
        currency = "IDR"
    }
    return paymentURL, txid, amountIDR, currency, nil
}

