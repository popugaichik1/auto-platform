package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	core_domain "messenger-service/internal/features/messenger/domain"
)

// ListMessages отдаёт страницу истории треда, от новых к старым.
func (r *Repository) ListMessages(
	ctx context.Context,
	conversationID uuid.UUID,
	page, limit int,
) ([]core_domain.Message, error) {
	op := "Messenger.Repo.ListMessages"

	opCtx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	const query = `
		SELECT id, conversation_id, sender_id, body, created_at
		FROM messengerservice.messages
		WHERE conversation_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(opCtx, query, conversationID, limit, offset)
	if err != nil {
		
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var messages []core_domain.Message
	for rows.Next() {
		var row messageRow
		if err := rows.Scan(
			&row.ID,
			&row.ConversationID,
			&row.SenderID,
			&row.Body,
			&row.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		messages = append(messages, row.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows: %w", op, err)
	}

	return messages, nil
}
