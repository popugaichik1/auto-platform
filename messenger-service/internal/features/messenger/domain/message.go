package core_domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
)

type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Body           string
	CreatedAt      time.Time
}

func NewMessage(conversationID, senderID uuid.UUID, body string) Message {
	return Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Body:           body,
		CreatedAt:      time.Now(),
	}
}

func (m Message) Validate() error {
	bodyLen := len([]rune(m.Body))
	if bodyLen < 1 || bodyLen > 4000 {
		return fmt.Errorf(
			"invalid `body` len: %d: %w",
			bodyLen,
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}
