package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
	listing_client "messenger-service/internal/clients/listing"
	core_domain "messenger-service/internal/features/messenger/domain"
)

func TestService_CreateOrGetConversation_SellerInitiates(t *testing.T) {
	listingID := uuid.New()
	sellerID := uuid.New()
	buyerID := uuid.New() // recipient, передан запросившим продавцом

	repo := &fakeRepo{
		createOrGetConversationFunc: func(ctx context.Context, conv core_domain.Conversation) (core_domain.Conversation, error) {
			if conv.SellerID != sellerID || conv.BuyerID != buyerID {
				t.Fatalf("CreateOrGetConversation() called with seller=%s buyer=%s, want seller=%s buyer=%s",
					conv.SellerID, conv.BuyerID, sellerID, buyerID)
			}
			return conv, nil
		},
	}
	listing := &fakeListingClient{
		getListingFunc: func(ctx context.Context, id uuid.UUID) (listing_client.Listing, error) {
			return listing_client.Listing{ID: listingID, UserID: sellerID}, nil
		},
	}
	svc := NewService(repo, listing, &fakePublisher{}, testLogger())

	conv, err := svc.CreateOrGetConversation(context.Background(), listingID, sellerID, buyerID)
	if err != nil {
		t.Fatalf("CreateOrGetConversation() error = %v", err)
	}
	if conv.SellerID != sellerID || conv.BuyerID != buyerID {
		t.Fatalf("CreateOrGetConversation() = %+v, want seller=%s buyer=%s", conv, sellerID, buyerID)
	}
}

func TestService_CreateOrGetConversation_BuyerInitiates(t *testing.T) {
	listingID := uuid.New()
	sellerID := uuid.New()
	buyerID := uuid.New() // requester — покупатель, пишет продавцу (recipient)

	repo := &fakeRepo{
		createOrGetConversationFunc: func(ctx context.Context, conv core_domain.Conversation) (core_domain.Conversation, error) {
			return conv, nil
		},
	}
	listing := &fakeListingClient{
		getListingFunc: func(ctx context.Context, id uuid.UUID) (listing_client.Listing, error) {
			return listing_client.Listing{ID: listingID, UserID: sellerID}, nil
		},
	}
	svc := NewService(repo, listing, &fakePublisher{}, testLogger())

	conv, err := svc.CreateOrGetConversation(context.Background(), listingID, buyerID, sellerID)
	if err != nil {
		t.Fatalf("CreateOrGetConversation() error = %v", err)
	}
	if conv.SellerID != sellerID || conv.BuyerID != buyerID {
		t.Fatalf("CreateOrGetConversation() = %+v, want seller=%s buyer=%s", conv, sellerID, buyerID)
	}
}

func TestService_CreateOrGetConversation_NeitherIsSeller(t *testing.T) {
	listingID := uuid.New()
	sellerID := uuid.New()
	stranger1 := uuid.New()
	stranger2 := uuid.New()

	listing := &fakeListingClient{
		getListingFunc: func(ctx context.Context, id uuid.UUID) (listing_client.Listing, error) {
			return listing_client.Listing{ID: listingID, UserID: sellerID}, nil
		},
	}
	svc := NewService(&fakeRepo{}, listing, &fakePublisher{}, testLogger())

	_, err := svc.CreateOrGetConversation(context.Background(), listingID, stranger1, stranger2)
	if !errors.Is(err, core_errors.ErrInvalidArgument) {
		t.Fatalf("CreateOrGetConversation() error = %v, want wrapped %v", err, core_errors.ErrInvalidArgument)
	}
}

func TestService_CreateOrGetConversation_ListingClientError(t *testing.T) {
	listing := &fakeListingClient{
		getListingFunc: func(ctx context.Context, id uuid.UUID) (listing_client.Listing, error) {
			return listing_client.Listing{}, core_errors.ErrNotFound
		},
	}
	svc := NewService(&fakeRepo{}, listing, &fakePublisher{}, testLogger())

	_, err := svc.CreateOrGetConversation(context.Background(), uuid.New(), uuid.New(), uuid.New())
	if !errors.Is(err, core_errors.ErrNotFound) {
		t.Fatalf("CreateOrGetConversation() error = %v, want wrapped %v", err, core_errors.ErrNotFound)
	}
}
