package entities

import (
	"time"

	"github.com/google/uuid"
)

type EmailTemplateEnum string

const (
	EmailTemplateVerificationEmail EmailTemplateEnum = "verification_email"
	EmailTemplateVerificationPhone EmailTemplateEnum = "verification_phone"
	EmailTemplatePasswordReset     EmailTemplateEnum = "password_reset"
)

type EmailTemplate struct {
	ID        uuid.UUID
	Name      EmailTemplateEnum
	Subject   string
	Body      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
