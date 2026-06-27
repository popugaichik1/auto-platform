package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	core_domain "user-service/internal/core/domain"
	core_errors "user-service/internal/core/errors"
	core_logger "user-service/internal/core/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func testLogger() *core_logger.Logger {
	return &core_logger.Logger{Logger: zap.NewNop()}
}

// fakeRepo — ручная fake-реализация Repo: каждый тест задаёт только то
// поле-функцию, которое ему нужно для конкретного сценария.
type fakeRepo struct {
	saveUserFunc    func(ctx context.Context, user core_domain.SaveUser) error
	getUserByIDFunc func(ctx context.Context, id uuid.UUID) (core_domain.User, error)
}

func (f *fakeRepo) SaveUser(ctx context.Context, user core_domain.SaveUser) error {
	return f.saveUserFunc(ctx, user)
}

func (f *fakeRepo) GetUserByID(ctx context.Context, id uuid.UUID) (core_domain.User, error) {
	return f.getUserByIDFunc(ctx, id)
}

func TestService_GetUser_Found(t *testing.T) {
	id := uuid.New()
	want := core_domain.User{ID: id, Version: 1, Username: "zahar", PhoneNumber: "+77001234567"}

	repo := &fakeRepo{
		getUserByIDFunc: func(ctx context.Context, gotID uuid.UUID) (core_domain.User, error) {
			if gotID != id {
				t.Fatalf("GetUserByID() called with %s, want %s", gotID, id)
			}
			return want, nil
		},
	}
	svc := NewService(repo, testLogger())

	got, err := svc.GetUser(context.Background(), id)
	if err != nil {
		t.Fatalf("GetUser() error = %v", err)
	}
	if got != want {
		t.Fatalf("GetUser() = %+v, want %+v", got, want)
	}
}

func TestService_GetUser_NotFound(t *testing.T) {
	id := uuid.New()
	repo := &fakeRepo{
		getUserByIDFunc: func(ctx context.Context, gotID uuid.UUID) (core_domain.User, error) {
			return core_domain.User{}, fmt.Errorf("repo: %w", core_errors.ErrNotFound)
		},
	}
	svc := NewService(repo, testLogger())

	_, err := svc.GetUser(context.Background(), id)
	if !errors.Is(err, core_errors.ErrNotFound) {
		t.Fatalf("GetUser() error = %v, want wrapped %v", err, core_errors.ErrNotFound)
	}
}
