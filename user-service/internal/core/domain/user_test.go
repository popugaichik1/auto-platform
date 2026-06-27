package core_domain

import (
	"errors"
	"strings"
	"testing"
	core_errors "user-service/internal/core/errors"

	"github.com/google/uuid"
)

func TestSaveUser_Validate(t *testing.T) {
	validUsername := "john_doe"
	validPhone := "+77001234567"

	tests := []struct {
		name        string
		username    string
		phoneNumber string
		wantErr     bool
	}{
		{
			name:        "valid user",
			username:    validUsername,
			phoneNumber: validPhone,
			wantErr:     false,
		},
		{
			name:        "empty username",
			username:    "",
			phoneNumber: validPhone,
			wantErr:     true,
		},
		{
			name:        "username too long",
			username:    strings.Repeat("a", 101),
			phoneNumber: validPhone,
			wantErr:     true,
		},
		{
			name:        "phone too short",
			username:    validUsername,
			phoneNumber: "+7700123",
			wantErr:     true,
		},
		{
			name:        "phone too long",
			username:    validUsername,
			phoneNumber: "+770012345678901",
			wantErr:     true,
		},
		{
			name:        "phone missing leading plus",
			username:    validUsername,
			phoneNumber: "77001234567",
			wantErr:     true,
		},
		{
			name:        "phone has non-digit characters",
			username:    validUsername,
			phoneNumber: "+7700123abc7",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewSaveUser(uuid.New(), tt.username, tt.phoneNumber)

			err := user.Validate()

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

func TestNewSaveUser_DefaultsVersionToOne(t *testing.T) {
	u := NewSaveUser(uuid.New(), "john_doe", "+77001234567")

	if u.Version != 1 {
		t.Fatalf("expected initial version 1, got %d", u.Version)
	}
}
