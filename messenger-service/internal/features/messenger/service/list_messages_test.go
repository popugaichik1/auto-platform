package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
	core_domain "messenger-service/internal/features/messenger/domain"
)

func TestService_ListMessages_Participant_OK(t *testing.T) {
	sellerID := uuid.New()
	conv := core_domain.NewConversation(uuid.New(), sellerID, uuid.New())
	want := []core_domain.Message{core_domain.NewMessage(conv.ID, sellerID, "hi")}

	repo := &fakeRepo{
		getConversationByIDFunc: func(ctx context.Context, id uuid.UUID) (core_domain.Conversation, error) {
			return conv, nil
		},
		listMessagesFunc: func(ctx context.Context, conversationID uuid.UUID, page, limit int) ([]core_domain.Message, error) {
			return want, nil
		},
	}
	svc := NewService(repo, &fakeListingClient{}, &fakePublisher{}, testLogger())

	got, err := svc.ListMessages(context.Background(), conv.ID, sellerID, 1, 20)
	if err != nil {
		t.Fatalf("ListMessages() error = %v", err)
	}
	if len(got) != 1 || got[0].Body != "hi" {
		t.Fatalf("ListMessages() = %+v, want %+v", got, want)
	}
}

func TestService_ListMessages_NonParticipant_Forbidden(t *testing.T) {
	conv := core_domain.NewConversation(uuid.New(), uuid.New(), uuid.New())
	stranger := uuid.New()

	repo := &fakeRepo{
		getConversationByIDFunc: func(ctx context.Context, id uuid.UUID) (core_domain.Conversation, error) {
			return conv, nil
		},
	}
	svc := NewService(repo, &fakeListingClient{}, &fakePublisher{}, testLogger())

	_, err := svc.ListMessages(context.Background(), conv.ID, stranger, 1, 20)
	if !errors.Is(err, core_errors.ErrForbidden) {
		t.Fatalf("ListMessages() error = %v, want wrapped %v", err, core_errors.ErrForbidden)
	}
}
