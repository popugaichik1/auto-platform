package repository

import (
	"context"
	"fmt"

	core_domain "messenger-service/internal/features/messenger/domain"
)

// CreateMessage сохраняет сообщение и продвигает last_message_at треда
// атомарно — в одной транзакции, чтобы между сохранением сообщения и
// обновлением метки активности треда не могло остаться неконсистентного
// состояния.
func (r *Repository) CreateMessage(
	ctx context.Context,
	msg core_domain.Message,
) (core_domain.Message, error) {
	op := "Messenger.Repo.CreateMessage"

	opCtx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(opCtx)
	if err != nil {
		return core_domain.Message{}, fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer func() { _ = tx.Rollback(opCtx) }() // no-op, если транзакция уже закоммичена

	const insertQuery = `
		INSERT INTO messengerservice.messages (
			id, conversation_id, sender_id, body
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, conversation_id, sender_id, body, created_at;
	`

	row := tx.QueryRow(
		opCtx,
		insertQuery,
		msg.ID,
		msg.ConversationID,
		msg.SenderID,
		msg.Body,
	)

	var result messageRow
	if err := row.Scan(
		&result.ID,
		&result.ConversationID,
		&result.SenderID,
		&result.Body,
		&result.CreatedAt,
	); err != nil {
		return core_domain.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	const touchQuery = `
		UPDATE messengerservice.conversations
		SET last_message_at = $2
		WHERE id = $1
	`
	if _, err := tx.Exec(opCtx, touchQuery, msg.ConversationID, result.CreatedAt); err != nil {
		return core_domain.Message{}, fmt.Errorf("%s: touch conversation: %w", op, err)
	}

	if err := tx.Commit(opCtx); err != nil {
		return core_domain.Message{}, fmt.Errorf("%s: commit: %w", op, err)
	}

	return result.toDomain(), nil
}
