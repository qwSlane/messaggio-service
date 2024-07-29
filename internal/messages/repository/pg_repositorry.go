package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"mesaggio-test/internal/messages"
	"mesaggio-test/internal/models"
)

type messageRepository struct {
	db *sqlx.DB
}

func NewMessagesRepository(db *sqlx.DB) messages.Repository {
	return &messageRepository{db: db}
}

func (m *messageRepository) SaveMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	var created models.Message
	if err := m.db.QueryRowxContext(ctx, createMessageQuery, &message.MessageID, &message.Content, &message.Status).Scan(
		&created.MessageID,
		&created.Content,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow")
	}

	return &created, nil
}

func (m *messageRepository) UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	return nil, nil
}
