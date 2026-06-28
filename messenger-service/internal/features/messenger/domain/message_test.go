package core_domain

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"

	core_errors "messenger-service/internal/core/errors"
)

func TestMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(m *Message)
		wantErr bool
	}{
		{
			name:    "valid message",
			mutate:  func(m *Message) {},
			wantErr: false,
		},
		{
			name:    "empty body",
			mutate:  func(m *Message) { m.Body = "" },
			wantErr: true,
		},
		{
			name:    "body too long",
			mutate:  func(m *Message) { m.Body = strings.Repeat("a", 4001) },
			wantErr: true,
		},
		{
			name:    "body at max length",
			mutate:  func(m *Message) { m.Body = strings.Repeat("a", 4000) },
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMessage(uuid.New(), uuid.New(), "hello")
			tt.mutate(&msg)

			err := msg.Validate()

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
