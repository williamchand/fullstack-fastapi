package repositories

import "context"

// WahaClient defines WhatsApp sender via WAHA API
type WahaClient interface {
    // SendText sends a text message to the given international phone number
    // Phone should be numeric; client formats to chatId `<number>@c.us`
    SendText(ctx context.Context, phone string, text string) error
}

