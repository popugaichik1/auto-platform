package repository

import (
	"context"
	"fmt"

	core_domain "messenger-service/internal/features/messenger/domain"
)

// CreateOrGetConversation — идемпотентный get-or-create: если тред с такими
// listing_id+buyer_id уже существует, возвращает его как есть (DO UPDATE —
// no-op, нужен только чтобы RETURNING всегда отдавал строку); иначе создаёт
// новый. seller_id фиксируется на первом создании и не переписывается даже
// если буквы запроса вдруг не совпадут — buyer не может стать продавцом.
func (r *Repository) CreateOrGetConversation(
	ctx context.Context,
	conv core_domain.Conversation,
) (core_domain.Conversation, error) {
	op := "Messenger.Repo.CreateOrGetConversation"

	opCtx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const query = `
		INSERT INTO messengerservice.conversations (
			id, listing_id, seller_id, buyer_id
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (listing_id, buyer_id) DO UPDATE SET
			last_message_at = messengerservice.conversations.last_message_at
		RETURNING id, listing_id, seller_id, buyer_id, created_at, last_message_at;
	`

	row := r.pool.QueryRow(
		opCtx,
		query,
		conv.ID,
		conv.ListingID,
		conv.SellerID,
		conv.BuyerID,
	)

	var result conversationRow
	if err := row.Scan(
		&result.ID,
		&result.ListingID,
		&result.SellerID,
		&result.BuyerID,
		&result.CreatedAt,
		&result.LastMessageAt,
	); err != nil {
		return core_domain.Conversation{}, fmt.Errorf("%s: %w", op, err)
	}

	return result.toDomain(), nil
}
