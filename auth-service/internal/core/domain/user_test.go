package core_domain

import (
	"errors"
	"strings"
	"testing"

	core_errors "github.com/zosinkin/social_network/internal/core/errors"
)

func TestAuthUser_Validate(t *testing.T) {
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
			// PasswordHash здесь не участвует в Validate() — пароль проверяется
			// отдельно функцией ValidatePassword до хеширования.
			user := NewAuthUser(tt.username, tt.phoneNumber, "irrelevant-hash")

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

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "valid password", password: "Sup3rSecret!", wantErr: false},
		{name: "too short", password: "Sh0rt!", wantErr: true},
		{name: "missing uppercase", password: "lowercase1!", wantErr: true},
		{name: "missing lowercase", password: "UPPERCASE1!", wantErr: true},
		{name: "missing digit or special char", password: "OnlyLetters", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ValidatePassword() expected error, got nil")
				}
				if !errors.Is(err, core_errors.ErrInvalidArgument) {
					t.Fatalf("ValidatePassword() error = %v, want wrapped %v", err, core_errors.ErrInvalidArgument)
				}
				return
			}

			if err != nil {
				t.Fatalf("ValidatePassword() unexpected error: %v", err)
			}
		})
	}
}

func TestNewAuthUser_GeneratesID(t *testing.T) {
	u1 := NewAuthUser("a", "+77001234567", "Sup3rSecret!")
	u2 := NewAuthUser("a", "+77001234567", "Sup3rSecret!")

	if u1.ID == u2.ID {
		t.Fatalf("expected distinct generated IDs, got the same: %v", u1.ID)
	}
}
