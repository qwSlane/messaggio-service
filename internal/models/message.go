package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Message struct {
	MessageID uuid.UUID `json:"messageID"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Status    bool      `json:"status"`
}

type CreatedResponse struct {
	MessageID uuid.UUID `json:"messageID"`
}
