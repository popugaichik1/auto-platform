package transport_http

import (
	"context"
	core_domain "user-service/internal/core/domain"
	core_logger "user-service/internal/core/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func testLogger() *core_logger.Logger {
	return &core_logger.Logger{Logger: zap.NewNop()}
}

// fakeService — ручная fake-реализация Service (см. transport.go): каждый
// тест задаёт только то поле-функцию, которое ему нужно.
type fakeService struct {
	saveUserFunc func(ctx context.Context, id uuid.UUID, username, phoneNumber string) error
	getUserFunc  func(ctx context.Context, id uuid.UUID) (core_domain.User, error)
}

func (f *fakeService) SaveUser(ctx context.Context, id uuid.UUID, username, phoneNumber string) error {
	return f.saveUserFunc(ctx, id, username, phoneNumber)
}

func (f *fakeService) GetUser(ctx context.Context, id uuid.UUID) (core_domain.User, error) {
	return f.getUserFunc(ctx, id)
}
