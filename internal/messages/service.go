package messages

import (
	"context"
	"mesaggio-test/internal/models"
)

type Service interface {
	ReceiveMessage(ctx context.Context, message *models.Message) error
	UpdateMessage(ctx context.Context) error
}
