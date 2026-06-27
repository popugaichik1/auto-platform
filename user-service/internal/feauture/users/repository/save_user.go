package repository

import (
	"context"
	core_domain "user-service/internal/core/domain"

	"go.uber.org/zap"
)


func (r *Repo) SaveUser(
	ctx context.Context,
	user core_domain.SaveUser,
) error {
	op := "User.Repo.Register"

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	// ON CONFLICT нужен потому, что Kafka гарантирует только at-least-once
	// доставку: то же событие user.registered может прийти повторно
	// (редеплой консьюмера, ребаланс группы и т.п.), и повтор должен быть
	// no-op, а не падением на UNIQUE(id) с потерей сообщения.
	query :=
	`
		INSERT INTO userservice.users (
			id, version, username, phone_number
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			username = EXCLUDED.username,
			phone_number = EXCLUDED.phone_number,
			updated_at = NOW()
		RETURNING id, version, username, phone_number;
	`

	row := r.pool.QueryRow(
		ctx, 
		query,
		user.ID,
		user.Version,
		user.Username,
		user.PhoneNumber,
	)

	var userModel UserModel
	if err := userModel.Scan(row); err != nil {
		r.log.Error("Failed to save user", zap.String("op", op), zap.Error(err))
		return err
	}

	r.log.Debug("User saved successfully:", zap.String("op", op))
	return nil
}