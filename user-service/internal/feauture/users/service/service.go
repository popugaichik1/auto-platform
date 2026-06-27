package service

import (
	"context"
	core_domain "user-service/internal/core/domain"
	core_logger "user-service/internal/core/logger"

	"github.com/google/uuid"
)


type Service struct {
	repo Repo
	log *core_logger.Logger
}

type Repo interface {
	SaveUser(
		ctx context.Context,
		user core_domain.SaveUser,
	) (error)

	GetUserByID(
		ctx context.Context,
		id uuid.UUID,
	) (core_domain.User, error)
}

func NewService(
	repo Repo,
	log *core_logger.Logger,
) *Service {
	return &Service{
		repo: repo,
		log: log,
	}
}
