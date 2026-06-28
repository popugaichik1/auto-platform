package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	core_logger "messenger-service/internal/core/logger"
	core_kafka "messenger-service/internal/core/transport/kafka"
	listing_client "messenger-service/internal/clients/listing"
	core_domain "messenger-service/internal/features/messenger/domain"
)

func testLogger() *core_logger.Logger {
	return &core_logger.Logger{Logger: zap.NewNop()}
}

// fakeRepo — ручная fake-реализация Repo: каждый тест задаёт только то
// поле-функцию, которое ему нужно для конкретного сценария.
type fakeRepo struct {
	createOrGetConversationFunc func(ctx context.Context, conv core_domain.Conversation) (core_domain.Conversation, error)
	getConversationByIDFunc     func(ctx context.Context, id uuid.UUID) (core_domain.Conversation, error)
	listConversationsFunc       func(ctx context.Context, userID uuid.UUID) ([]core_domain.Conversation, error)
	createMessageFunc           func(ctx context.Context, msg core_domain.Message) (core_domain.Message, error)
	listMessagesFunc            func(ctx context.Context, conversationID uuid.UUID, page, limit int) ([]core_domain.Message, error)
}

func (f *fakeRepo) CreateOrGetConversation(ctx context.Context, conv core_domain.Conversation) (core_domain.Conversation, error) {
	return f.createOrGetConversationFunc(ctx, conv)
}

func (f *fakeRepo) GetConversationByID(ctx context.Context, id uuid.UUID) (core_domain.Conversation, error) {
	return f.getConversationByIDFunc(ctx, id)
}

func (f *fakeRepo) ListConversations(ctx context.Context, userID uuid.UUID) ([]core_domain.Conversation, error) {
	return f.listConversationsFunc(ctx, userID)
}

func (f *fakeRepo) CreateMessage(ctx context.Context, msg core_domain.Message) (core_domain.Message, error) {
	return f.createMessageFunc(ctx, msg)
}

func (f *fakeRepo) ListMessages(ctx context.Context, conversationID uuid.UUID, page, limit int) ([]core_domain.Message, error) {
	return f.listMessagesFunc(ctx, conversationID, page, limit)
}

// fakeListingClient — fake источника "кто продавец".
type fakeListingClient struct {
	getListingFunc func(ctx context.Context, id uuid.UUID) (listing_client.Listing, error)
}

func (f *fakeListingClient) GetListing(ctx context.Context, id uuid.UUID) (listing_client.Listing, error) {
	return f.getListingFunc(ctx, id)
}

// fakePublisher — по умолчанию (без заданного publishFunc) считает
// публикацию успешной, чтобы тестам, не проверяющим Kafka, не нужно было
// задавать функцию вовсе.
type fakePublisher struct {
	publishFunc func(ctx context.Context, message core_kafka.Message) error
}

func (f *fakePublisher) Publish(ctx context.Context, message core_kafka.Message) error {
	if f.publishFunc == nil {
		return nil
	}
	return f.publishFunc(ctx, message)
}
