package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	core_domain "messenger-service/internal/features/messenger/domain"
)

// ListConversations возвращает треды, где userID — продавец или покупатель,
// отсортированные по последней активности.
func (r *Repository) ListConversations(
	ctx context.Context,
	userID uuid.UUID,
) ([]core_domain.Conversation, error) {
	op := "Messenger.Repo.ListConversations"

	opCtx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	const query = `
		SELECT id, listing_id, seller_id, buyer_id, created_at, last_message_at
		FROM messengerservice.conversations
		WHERE seller_id = $1 OR buyer_id = $1
		ORDER BY last_message_at DESC
	`

	rows, err := r.pool.Query(opCtx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var conversations []core_domain.Conversation
	for rows.Next() {
		var row conversationRow
		if err := rows.Scan(
			&row.ID,
			&row.ListingID,
			&row.SellerID,
			&row.BuyerID,
			&row.CreatedAt,
			&row.LastMessageAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		conversations = append(conversations, row.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows: %w", op, err)
	}

	return conversations, nil
}
