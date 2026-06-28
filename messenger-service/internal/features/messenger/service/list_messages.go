package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "messenger-service/internal/core/errors"
	core_domain "messenger-service/internal/features/messenger/domain"
)

func (s *Service) ListMessages(
	ctx context.Context,
	conversationID uuid.UUID,
	requesterID uuid.UUID,
	page, limit int,
) ([]core_domain.Message, error) {
	op := "Messenger.Service.ListMessages"

	conv, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !conv.IsParticipant(requesterID) {
		return nil, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
	}

	messages, err := s.repo.ListMessages(ctx, conversationID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, nil
}
