package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_domain "github.com/zosinkin/social_network/internal/core/domain"
	core_errors "github.com/zosinkin/social_network/internal/core/errors"
	core_postgres_pool "github.com/zosinkin/social_network/internal/core/repository/postgres/pool"
	"go.uber.org/zap"
)


func (r *RefreshTokenRepo) GetRefreshToken(
	ctx context.Context,
	tokenString string,
) (*core_domain.RefreshToken, error) {
	op := "AuthService.Repo.GetRefreshToken"

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked
		FROM authservice.refresh_tokens
		WHERE token = $1
	`
	row := r.pool.QueryRow(ctx, query, tokenString)

	var RefreshToken RefreshTokenModel
	if err := RefreshToken.TokenScan(row); err != nil {
		r.log.Error("Token scan error:", zap.String("op", op), zap.Error(err))
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return &core_domain.RefreshToken{}, fmt.Errorf(
				"%v: %v",
				op,
				core_errors.ErrNotFound,
			)
		}
	}
	
	RefreshTokenDomain := modelTokenDomain(RefreshToken)

	return &RefreshTokenDomain, nil
	
}


func (r *RefreshTokenRepo) CreateRefreshToken(
	ctx context.Context,
	token core_domain.RefreshToken,
) (core_domain.RefreshToken, error) {
	op := "AuthService.Repo.CreateRefreshToken"

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO authservice.refresh_tokens (
			id, user_id, token, expires_at, created_at, revoked
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, token, expires_at, created_at, revoked
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
		token.Revoked,
	)
	var refreshTokenModel RefreshTokenModel
	if err := refreshTokenModel.TokenScan(row); err != nil {
		r.log.Error("Token scan error:", zap.String("op", op), zap.Error(err))
		return core_domain.RefreshToken{}, fmt.Errorf("%s: %w", op, err)
	}
	

	return token, nil
}


func (r *RefreshTokenRepo) RevokeRefreshToken(
	ctx context.Context,
	tokenString string,
) error {
	op := "AuthServie.Repo.RevokeRefreshToken"
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE authservice.refresh_tokens
		SET revoked = true
		WHERE token = $1
	`

	_, err := r.pool.Exec(ctx, query, tokenString)
	if err != nil {
		r.log.Error("Revoke refresh token error:", zap.String("op", op), zap.Error(err))
		return fmt.Errorf("%v: %w", op, err)
	}
	return nil
	
}