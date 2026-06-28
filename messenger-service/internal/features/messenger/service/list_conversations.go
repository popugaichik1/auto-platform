package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_domain "messenger-service/internal/features/messenger/domain"
)

func (s *Service) ListConversations(
	ctx context.Context,
	userID uuid.UUID,
) ([]core_domain.Conversation, error) {
	op := "Messenger.Service.ListConversations"

	conversations, err := s.repo.ListConversations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return conversations, nil
}
