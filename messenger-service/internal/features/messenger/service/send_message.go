package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	core_errors "messenger-service/internal/core/errors"
	core_kafka "messenger-service/internal/core/transport/kafka"
	core_domain "messenger-service/internal/features/messenger/domain"
	transport_kafka "messenger-service/internal/features/messenger/transport/kafka"
)

func (s *Service) SendMessage(
	ctx context.Context,
	conversationID uuid.UUID,
	senderID uuid.UUID,
	body string,
) (core_domain.Message, error) {
	op := "Messenger.Service.SendMessage"

	conv, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return core_domain.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	if !conv.IsParticipant(senderID) {
		return core_domain.Message{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
	}

	msg := core_domain.NewMessage(conversationID, senderID, body)
	if err := msg.Validate(); err != nil {
		return core_domain.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	saved, err := s.repo.CreateMessage(ctx, msg)
	if err != nil {
		return core_domain.Message{}, fmt.Errorf("%s: %w", op, err)
	}

	recipientID := conv.SellerID
	if senderID == conv.SellerID {
		recipientID = conv.BuyerID
	}

	event := transport_kafka.MessageSentEvent{
		MessageID:      saved.ID,
		ConversationID: saved.ConversationID,
		SenderID:       saved.SenderID,
		RecipientID:    recipientID,
		Body:           saved.Body,
		CreatedAt:      saved.CreatedAt,
	}

	_ = s.publisher.Publish(ctx, core_kafka.NewMessage(
		core_kafka.TopicMessageSent,
		saved.ID.String(),
		event,
	))

	return saved, nil
}
