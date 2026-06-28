package core_domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
)

type Conversation struct {
	ID            uuid.UUID
	ListingID     uuid.UUID
	SellerID      uuid.UUID
	BuyerID       uuid.UUID
	CreatedAt     time.Time
	LastMessageAt time.Time
}

func NewConversation(listingID, sellerID, buyerID uuid.UUID) Conversation {
	now := time.Now()
	return Conversation{
		ID:            uuid.New(),
		ListingID:     listingID,
		SellerID:      sellerID,
		BuyerID:       buyerID,
		CreatedAt:     now,
		LastMessageAt: now,
	}
}

func (c Conversation) Validate() error {
	if c.SellerID == c.BuyerID {
		return fmt.Errorf("seller and buyer cannot be the same user: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

// IsParticipant сообщает, является ли пользователь одной из двух сторон треда.
func (c Conversation) IsParticipant(userID uuid.UUID) bool {
	return userID == c.SellerID || userID == c.BuyerID
}
