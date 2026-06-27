package service

import (
	"context"
	"fmt"
	core_domain "user-service/internal/core/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (core_domain.User, error) {
	op := "User.Service.GetUser"

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		s.log.Error("get user error", zap.String("op", op), zap.Error(err))
		return core_domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
