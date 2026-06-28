package transport_kafka

import (
	"time"

	"github.com/google/uuid"
)

// MessageSentEvent публикуется в TopicMessageSent после сохранения
// сообщения. Используется и продюсером (service.SendMessage), и
// консьюмером фан-аута — обе стороны находятся в одном бинарнике, поэтому
// делят один Go-тип вместо неявного JSON-контракта между сервисами.
type MessageSentEvent struct {
	MessageID      uuid.UUID `json:"message_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	RecipientID    uuid.UUID `json:"recipient_id"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
}
