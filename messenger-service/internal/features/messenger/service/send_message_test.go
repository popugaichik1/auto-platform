package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
	core_kafka "messenger-service/internal/core/transport/kafka"
	core_domain "messenger-service/internal/features/messenger/domain"
	transport_kafka "messenger-service/internal/features/messenger/transport/kafka"
)

func TestService_SendMessage_BySeller_RecipientIsBuyer(t *testing.T) {
	sellerID := uuid.New()
	buyerID := uuid.New()
	conv := core_domain.NewConversation(uuid.New(), sellerID, buyerID)

	var publishedEvent transport_kafka.MessageSentEvent
	repo := &fakeRepo{
		getConversationByIDFunc: func(ctx context.Context, id uuid.UUID) (core_domain.Conversation, error) {
			return conv, nil
		},
		createMessageFunc: func(ctx context.Context, msg core_domain.Message) (core_domain.Message, error) {
			return msg, nil
		},
	}
	publisher := &fakePublisher{
		publishFunc: func(ctx context.Context, message core_kafka.Message) error {
			publishedEvent = message.Payload.(transport_kafka.MessageSentEvent)
			return nil
		},
	}
	svc := NewService(repo, &fakeListingClient{}, publisher, testLogger())

	msg, err := svc.SendMessage(context.Background(), conv.ID, sellerID, "hello")
	if err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}
	if msg.SenderID != sellerID {
		t.Fatalf("SendMessage() sender = %s, want %s", msg.SenderID, sellerID)
	}
	if publishedEvent.RecipientID != buyerID {
		t.Fatalf("published event recipient = %s, want buyer %s", publishedEvent.RecipientID, buyerID)
	}
}

func TestService_SendMessage_ByBuyer_RecipientIsSeller(t *testing.T) {
	sellerID := uuid.New()
	buyerID := uuid.New()
	conv := core_domain.NewConversation(uuid.New(), sellerID, buyerID)

	var publishedEvent transport_kafka.MessageSentEvent
	repo := &fakeRepo{
		getConversationByIDFunc: func(ctx context.Context, id uuid.UUID) (core_domain.Conversation, error) {
			return conv, nil
		},
		createMessageFunc: func(ctx context.Context, msg core_domain.Message) (core_domain.Message, error) {
			return msg, nil
		},
	}
	publisher := &fakePublisher{
		publishFunc: func(ctx context.Context, message core_kafka.Message) error {
			publishedEvent = message.Payload.(transport_kafka.MessageSentEvent)
			return nil
		},
	}
	svc := NewService(repo, &fakeListingClient{}, publisher, testLogger())

	_, err := svc.SendMessage(context.Background(), conv.ID, buyerID, "hello")
	if err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}
	if publishedEvent.RecipientID != sellerID {
		t.Fatalf("published event recipient = %s, want seller %s", publishedEvent.RecipientID, sellerID)
	}
}

func TestService_SendMessage_NonParticipant_Forbidden(t *testing.T) {
	conv := core_domain.NewConversation(uuid.New(), uuid.New(), uuid.New())
	stranger := uuid.New()

	repo := &fakeRepo{
		getConversationByIDFunc: func(ctx context.Context, id uuid.UUID) (core_domain.Conversation, error) {
			return conv, nil
		},
	}
	svc := NewService(repo, &fakeListingClient{}, &fakePublisher{}, testLogger())

	_, err := svc.SendMessage(context.Background(), conv.ID, stranger, "hello")
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("SendMessage() error = %v, want wrapped %v", err, core_errors.ErrForbidden)
	}
}

func TestService_SendMessage_InvalidBody(t *testing.T) {
	sellerID := uuid.New()
	conv := core_domain.NewConversation(uuid.New(), sellerID, uuid.New())

	repo := &fakeRepo{
		getConversationByIDFunc: func(ctx context.Context, id uuid.UUID) (core_domain.Conversation, error) {
			return conv, nil
		},
	}
	svc := NewService(repo, &fakeListingClient{}, &fakePublisher{}, testLogger())

	_, err := svc.SendMessage(context.Background(), conv.ID, sellerID, "")
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("SendMessage() error = %v, want wrapped %v", err, core_errors.ErrInvalidArgument)
	}
}

func TestService_SendMessage_PublishFailure_StillReturnsSavedMessage(t *testing.T) {
	sellerID := uuid.New()
	conv := core_domain.NewConversation(uuid.New(), sellerID, uuid.New())

	repo := &fakeRepo{
		getConversationByIDFunc: func(ctx context.Context, id uuid.UUID) (core_domain.Conversation, error) {
			return conv, nil
		},
		createMessageFunc: func(ctx context.Context, msg core_domain.Message) (core_domain.Message, error) {
			return msg, nil
		},
	}
	publisher := &fakePublisher{
		publishFunc: func(ctx context.Context, message core_kafka.Message) error {
			return errors.New("kafka is down")
		},
	}
	svc := NewService(repo, &fakeListingClient{}, publisher, testLogger())

	msg, err := svc.SendMessage(context.Background(), conv.ID, sellerID, "hello")
	if err != nil {
		t.Fatalf("SendMessage() error = %v, want nil (publish failure must not fail the send)", err)
	}
	if msg.Body != "hello" {
		t.Fatalf("SendMessage() = %+v, want saved message", msg)
	}
}
