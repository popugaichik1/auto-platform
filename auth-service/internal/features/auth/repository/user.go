package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	core_domain "github.com/zosinkin/social_network/internal/core/domain"
	core_postgres_pool "github.com/zosinkin/social_network/internal/core/repository/postgres/pool"
	"go.uber.org/zap"
)


func (r *UserRepo) RegisterUser(
	ctx context.Context,
	user core_domain.AuthUser,
) (core_domain.AuthUser, error) {
	op := "AuthService.Repo.RegisterUser"

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO authservice.users (
			id, username, phone_number, password_hash
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, phone_number;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.ID,
		user.Username,
		user.PhoneNumber,
		user.PasswordHash,
	)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		r.log.Error("UserModel scan error:", zap.String("op", op), zap.Error(err))
		return core_domain.AuthUser{}, fmt.Errorf("%s: %w", op, err)
	}

	userDomain := modelToUserDomain(userModel)

	return userDomain, nil
}


func (r *UserRepo) GetUserByPhoneNumber(
	ctx context.Context,
	phoneNumber string,
) (core_domain.AuthUser, error) {
	op := "AuthService.Repo.GetUserByPhoneNumber"

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := 	`
					SELECT id, username, phone_number, password_hash
					FROM authservice.users
					WHERE phone_number = $1
				`
	row := r.pool.QueryRow(
		ctx, 
		query, 
		phoneNumber,
	)

	var userModel UserModel
	if err := userModel.ScanWithPassword(row); err != nil {
		r.log.Error("Scan with password error:", zap.String("op", op), zap.Error(err))
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return core_domain.AuthUser{}, fmt.Errorf("%v: %w", op, err)
		}
		return core_domain.AuthUser{}, fmt.Errorf("%v: %w", op, err) 
	}
	userDomain := modelToUserDomain(userModel)

	return userDomain, nil
}


func (r *UserRepo) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (core_domain.AuthUser, error) {
	op := "AuthService.Repo.GetUserByID"

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, username, phone_number
		FROM authservice.users
		WHERE id = $1
	`

	row := r.pool.QueryRow(
		ctx, 
		query,
		id,
	)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		r.log.Error("user model scan error:", zap.String("op", op), zap.Error(err))
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return core_domain.AuthUser{}, fmt.Errorf("%v: %v", op, err)
		}
		return core_domain.AuthUser{}, fmt.Errorf("%v: %v", op, err)
	} 
	userDomain := modelToUserDomain(userModel)

	return userDomain, nil
}