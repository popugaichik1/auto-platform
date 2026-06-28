package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
	core_postgres_pool "messenger-service/internal/core/repository/postgres/pool"
	core_domain "messenger-service/internal/features/messenger/domain"
)

func (r *Repository) GetConversationByID(
	ctx context.Context,
	id uuid.UUID,
) (core_domain.Conversation, error) {
	op := "Messenger.Repo.GetConversationByID"

	opCtx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const query = `
		SELECT id, listing_id, seller_id, buyer_id, created_at, last_message_at
		FROM messengerservice.conversations
		WHERE id = $1
	`

	row := r.pool.QueryRow(opCtx, query, id)

	var result conversationRow
	if err := row.Scan(
		&result.ID,
		&result.ListingID,
		&result.SellerID,
		&result.BuyerID,
		&result.CreatedAt,
		&result.LastMessageAt,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return core_domain.Conversation{}, fmt.Errorf("%s: %w", op, core_errors.ErrNotFound)
		}
		return core_domain.Conversation{}, fmt.Errorf("%s: %w", op, err)
	}

	return result.toDomain(), nil
}
