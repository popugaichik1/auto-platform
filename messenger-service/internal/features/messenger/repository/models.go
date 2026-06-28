package repository

import (
	"time"

	"github.com/google/uuid"

	core_domain "messenger-service/internal/features/messenger/domain"
)

type conversationRow struct {
	ID            uuid.UUID
	ListingID     uuid.UUID
	SellerID      uuid.UUID
	BuyerID       uuid.UUID
	CreatedAt     time.Time
	LastMessageAt time.Time
}

func (r conversationRow) toDomain() core_domain.Conversation {
	return core_domain.Conversation{
		ID:            r.ID,
		ListingID:     r.ListingID,
		SellerID:      r.SellerID,
		BuyerID:       r.BuyerID,
		CreatedAt:     r.CreatedAt,
		LastMessageAt: r.LastMessageAt,
	}
}

type messageRow struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Body           string
	CreatedAt      time.Time
}

func (r messageRow) toDomain() core_domain.Message {
	return core_domain.Message{
		ID:             r.ID,
		ConversationID: r.ConversationID,
		SenderID:       r.SenderID,
		Body:           r.Body,
		CreatedAt:      r.CreatedAt,
	}
}
