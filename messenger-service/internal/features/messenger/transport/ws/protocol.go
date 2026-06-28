package ws

import "github.com/google/uuid"

// ClientFrame — то, что клиент шлёт серверу по уже установленному соединению.
type ClientFrame struct {
	Type           string    `json:"type"` // "send_message"
	ConversationID uuid.UUID `json:"conversation_id"`
	Body           string    `json:"body"`
}

// ServerFrame — то, что сервер шлёт клиенту: подтверждение отправки,
// входящее сообщение от собеседника, либо ошибка.
type ServerFrame struct {
	Type    string `json:"type"` // "message_sent" | "message" | "error"
	Payload any    `json:"payload"`
}
