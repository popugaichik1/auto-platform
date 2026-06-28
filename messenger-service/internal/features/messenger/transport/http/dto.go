package transport_http

import (
	"time"

	"github.com/google/uuid"

	core_domain "messenger-service/internal/features/messenger/domain"
)

type CreateConversationRequest struct {
	ListingID   uuid.UUID `json:"listing_id" binding:"required"`
	RecipientID uuid.UUID `json:"recipient_id" binding:"required"`
}

type ConversationResponse struct {
	ID            uuid.UUID `json:"id"`
	ListingID     uuid.UUID `json:"listing_id"`
	SellerID      uuid.UUID `json:"seller_id"`
	BuyerID       uuid.UUID `json:"buyer_id"`
	CreatedAt     time.Time `json:"created_at"`
	LastMessageAt time.Time `json:"last_message_at"`
}

func toConversationResponse(c core_domain.Conversation) ConversationResponse {
	return ConversationResponse{
		ID:            c.ID,
		ListingID:     c.ListingID,
		SellerID:      c.SellerID,
		BuyerID:       c.BuyerID,
		CreatedAt:     c.CreatedAt,
		LastMessageAt: c.LastMessageAt,
	}
}

type MessageResponse struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
}

func toMessageResponse(m core_domain.Message) MessageResponse {
	return MessageResponse{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		Body:           m.Body,
		CreatedAt:      m.CreatedAt,
	}
}
