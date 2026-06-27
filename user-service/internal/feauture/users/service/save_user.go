package service

import (
	"context"
	core_domain "user-service/internal/core/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)


func (s *Service) SaveUser(
	ctx context.Context,
	id uuid.UUID,
	username string,
	phoneNumber string,
) (error) {
	op := "User.Service.SaveUser"

	user := core_domain.NewSaveUser(
		id,
		username,
		phoneNumber,
	)

	if err := user.Validate(); err != nil {
		s.log.Error("validation error", zap.String("op", op), zap.Error(err))
		return err
	}

	if err := s.repo.SaveUser(ctx, user); err != nil {
		s.log.Error("save user error", zap.String("op", op), zap.Error(err))
		return err
	}

	return nil
}