package repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "user-service/internal/core/domain"
	core_errors "user-service/internal/core/errors"
	core_postgres_pool "user-service/internal/core/repository/postgres/pool"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (r *Repo) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (core_domain.User, error) {
	op := "User.Repo.GetUserByID"

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, username, phone_number
		FROM userservice.users
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return core_domain.User{}, fmt.Errorf("%s: %w", op, core_errors.ErrNotFound)
		}
		r.log.Error("get user by id error", zap.String("op", op), zap.Error(err))
		return core_domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return core_domain.User{
		ID:          userModel.ID,
		Version:     userModel.Version,
		Username:    userModel.Username,
		PhoneNumber: userModel.PhoneNumber,
	}, nil
}
