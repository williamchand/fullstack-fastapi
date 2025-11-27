package repositories

import "github.com/williamchand/fullstack-fastapi/backend-go/internal/domain/entities"

type Sender interface {
	Send(msg entities.Message) error
}
