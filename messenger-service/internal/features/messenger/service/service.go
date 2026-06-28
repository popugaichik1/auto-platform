package service

import (
	"context"

	"github.com/google/uuid"

	core_logger "messenger-service/internal/core/logger"
	core_kafka "messenger-service/internal/core/transport/kafka"
	listing_client "messenger-service/internal/clients/listing"
	core_domain "messenger-service/internal/features/messenger/domain"
)

type Service struct {
	repo      Repo
	listing   ListingClient
	publisher EventPublisher
}

type Repo interface {
	CreateOrGetConversation(ctx context.Context, conv core_domain.Conversation) (core_domain.Conversation, error)
	GetConversationByID(ctx context.Context, id uuid.UUID) (core_domain.Conversation, error)
	ListConversations(ctx context.Context, userID uuid.UUID) ([]core_domain.Conversation, error)
	CreateMessage(ctx context.Context, msg core_domain.Message) (core_domain.Message, error)
	ListMessages(ctx context.Context, conversationID uuid.UUID, page, limit int) ([]core_domain.Message, error)
}

// ListingClient — то немногое, что сервису нужно знать о другом сервисе:
// кто владелец (продавец) объявления.
type ListingClient interface {
	GetListing(ctx context.Context, id uuid.UUID) (listing_client.Listing, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, message core_kafka.Message) error
}

func NewService(repo Repo, listing ListingClient, publisher EventPublisher, log *core_logger.Logger) *Service {
	return &Service{
		repo:      repo,
		listing:   listing,
		publisher: publisher,
	}
}
