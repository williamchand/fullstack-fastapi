package entities

import (
	"time"

	"github.com/google/uuid"
)

type VerificationType string
type VerificationPurpose string

const (
	VerificationTypeEmail             VerificationType = "email"
	VerificationTypePhone             VerificationType = "phone"
	VerificationTypePasswordReset     VerificationType = "password_reset"
	VerificationTypePhoneRegistration VerificationType = "phone_registration"
)

const (
	VerificationPurposeEmailVerification VerificationPurpose = "email_verification"
	VerificationPurposeAddEmail          VerificationPurpose = "add_email"
	VerificationPurposeAddPhone          VerificationPurpose = "add_phone"
	VerificationPurposePhoneOTP          VerificationPurpose = "phone_otp"
	VerificationPurposePasswordReset     VerificationPurpose = "password_reset"
)

type VerificationCode struct {
	ID            uuid.UUID
	UserID        *uuid.UUID
	Code          string
	Type          VerificationType
	ExtraMetadata map[string]any
	CreatedAt     time.Time
	ExpiresAt     time.Time
	UsedAt        *time.Time
}
