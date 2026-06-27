package core_domain

import (
	"fmt"
	"regexp"
	"github.com/google/uuid"
	core_errors "github.com/zosinkin/social_network/internal/core/errors"
)


type AuthUser struct {
	ID  			uuid.UUID
	Username 		string
	PhoneNumber 	string
	PasswordHash    string
}



// CreateUser создаёт нового пользователя с автоматически сгенерированными
// ID (UUID v4) и начальной версией 1.
func NewAuthUser(
	username string,
	phoneNumber string,
	passwordHash string,
) AuthUser {
	var (
		id      = uuid.New()
	)

	return AuthUser{
		id,
		username,
		phoneNumber,
		passwordHash,
	}
}

func (u AuthUser) Validate() error {
	usernameLen := len([]rune(u.Username))
	if usernameLen < 1 || usernameLen > 100 {
		return fmt.Errorf(
			"inavlid `username` len: %d: %w",
			usernameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	phoneNumberLen := len([]rune(u.PhoneNumber))
	if phoneNumberLen < 10 || phoneNumberLen > 15 {
		return fmt.Errorf(
			"invalid `phoneNumber` len: %d: %w",
			phoneNumberLen,
			core_errors.ErrInvalidArgument,
		)
	}

	phonere := regexp.MustCompile(`^\+[0-9]+$`)

	if !phonere.MatchString(u.PhoneNumber) {
		return fmt.Errorf(
			"invalid `PhoneNumber` format: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

// ValidatePassword проверяет сырой (нехешированный) пароль на соответствие
// парольной политике. Должна вызываться до HashPassword — после хеширования
// эти проверки (наличие цифр/регистра) теряют смысл, так как bcrypt-хеш
// уже не несёт семантики исходного пароля.
func ValidatePassword(password string) error {
	passlen := len([]rune(password))
	if passlen < 8 {
		return fmt.Errorf(
			"`password` must be minimum 8 symbols: %d: %w",
			passlen,
			core_errors.ErrInvalidArgument,
		)
	}

	hasUpper := regexp.MustCompile(`[A-Z]`)
	hasLower := regexp.MustCompile(`[a-z]`)
	hasDigitOrSpecial := regexp.MustCompile(`[0-9\W]`)

	if !hasUpper.MatchString(password) ||
		!hasLower.MatchString(password) ||
		!hasDigitOrSpecial.MatchString(password) {
		return fmt.Errorf(
			"invalid `password` format: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}