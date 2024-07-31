package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// Message model
type Message struct {
	MessageID uuid.UUID `json:"messageID" db:"message_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	Status    bool      `json:"status" db:"status"`
}

type CreatedResponse struct {
	MessageID uuid.UUID `json:"messageID"`
}

type Request struct {
	Content string `json:"content"`
}

type Statistics struct {
	AllMessagesCount       int64 `json:"allMessages" db:"all_messages_count"`
	ProcessedMessagesCount int64 `json:"processedMessages" db:"processed_messages_count"`

	LastProcessedID uuid.UUID `json:"lastProcessed" db:"last_processed_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
