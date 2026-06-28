package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
	core_domain "messenger-service/internal/features/messenger/domain"
)

// CreateOrGetConversation узнаёт продавца у listing-service и определяет,
// кто из двух сторон запроса — продавец, а кто покупатель, прежде чем
// создать (или вернуть существующий) тред.
func (s *Service) CreateOrGetConversation(
	ctx context.Context,
	listingID uuid.UUID,
	requesterID uuid.UUID,
	recipientID uuid.UUID,
) (core_domain.Conversation, error) {
	op := "Messenger.Service.CreateOrGetConversation"

	listing, err := s.listing.GetListing(ctx, listingID)
	if err != nil {
		return core_domain.Conversation{}, fmt.Errorf("%s: %w", op, err)
	}
	sellerID := listing.UserID

	var buyerID uuid.UUID
	switch {
	case requesterID == sellerID:
		// продавец пишет покупателю
		buyerID = recipientID
	case recipientID == sellerID:
		// покупатель пишет продавцу
		buyerID = requesterID
	default:
		// ни одна из сторон запроса не владеет объявлением — бессмысленный тред
		return core_domain.Conversation{}, fmt.Errorf(
			"%s: recipient is not the listing owner: %w",
			op, core_errors.ErrInvalidArgument,
		)
	}

	conv := core_domain.NewConversation(listingID, sellerID, buyerID)
	if err := conv.Validate(); err != nil {
		return core_domain.Conversation{}, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.repo.CreateOrGetConversation(ctx, conv)
	if err != nil {
		return core_domain.Conversation{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
