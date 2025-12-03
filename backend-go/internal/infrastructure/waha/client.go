package waha

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    "time"
)

// Client implements WAHA WhatsApp API integration
type Client struct {
    url       string
    apiKey    string
    session   string
    httpClient *http.Client
}

func New(url, apiKey, session string) *Client {
    return &Client{
        url:       strings.TrimRight(url, "/"),
        apiKey:    apiKey,
        session:   session,
        httpClient: &http.Client{Timeout: 10 * time.Second},
    }
}

type sendTextPayload struct {
    Session string `json:"session"`
    ChatID  string `json:"chatId"`
    Text    string `json:"text"`
}

func (c *Client) SendText(ctx context.Context, phone string, text string) error {
    chatID := formatChatID(phone)
    payload := sendTextPayload{Session: c.session, ChatID: chatID, Text: text}
    b, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("waha: marshal payload failed: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+"/api/sendText", bytes.NewReader(b))
    if err != nil {
        return fmt.Errorf("waha: request build failed: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    if c.apiKey != "" {
        req.Header.Set("X-Api-Key", c.apiKey)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("waha: request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 {
        return fmt.Errorf("waha: sendText failed with status %d", resp.StatusCode)
    }
    return nil
}

// formatChatID converts an international phone to WAHA chatId: remove '+' and non-digits, append '@c.us'
func formatChatID(phone string) string {
    // strip spaces, dashes, parentheses
    cleaned := strings.Map(func(r rune) rune {
        if r >= '0' && r <= '9' {
            return r
        }
        return -1
    }, phone)
    // ensure no leading '+' remains; cleaned already removed non-digits
    return cleaned + "@c.us"
}

