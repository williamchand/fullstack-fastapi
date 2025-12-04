package repositories

import (
	"context"

	"github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"
)

type EmailTemplateRepository interface {
	TxProvider[EmailTemplateRepository]

	GetByName(ctx context.Context, name entities.EmailTemplateEnum) (*entities.EmailTemplate, error)
}
