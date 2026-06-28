package core_domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
)

func TestConversation_Validate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(c *Conversation)
		wantErr bool
	}{
		{
			name:    "valid conversation",
			mutate:  func(c *Conversation) {},
			wantErr: false,
		},
		{
			name: "seller and buyer are the same user",
			mutate: func(c *Conversation) {
				c.BuyerID = c.SellerID
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := NewConversation(uuid.New(), uuid.New(), uuid.New())
			tt.mutate(&conv)

			err := conv.Validate()

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Validate() expected error, got nil")
				}
				if !errors.Is(err, core_errors.ErrInvalidArgument) {
					t.Fatalf("Validate() error = %v, want wrapped %v", err, core_errors.ErrInvalidArgument)
				}
				return
			}

			if err != nil {
				t.Fatalf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestConversation_IsParticipant(t *testing.T) {
	sellerID := uuid.New()
	buyerID := uuid.New()
	stranger := uuid.New()

	conv := NewConversation(uuid.New(), sellerID, buyerID)

	if !conv.IsParticipant(sellerID) {
		t.Fatalf("IsParticipant(sellerID) = false, want true")
	}
	if !conv.IsParticipant(buyerID) {
		t.Fatalf("IsParticipant(buyerID) = false, want true")
	}
	if conv.IsParticipant(stranger) {
		t.Fatalf("IsParticipant(stranger) = true, want false")
	}
}
