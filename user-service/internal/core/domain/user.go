package core_domain

import (
	"fmt"
	"regexp"
	core_errors "user-service/internal/core/errors"

	"github.com/google/uuid"
)


type SaveUser struct {
	ID  			uuid.UUID
	Version 		int
	Username 		string
	PhoneNumber 	string
}

type User struct {
	ID  			uuid.UUID
	Version 		int
	Username 		string
	PhoneNumber 	string
}


func NewUser(
	ID 			uuid.UUID,
	Version 	int,
	Username 	string,
	PhoneNumber string,
) *User {
	return &User{
		ID: ID,
		Version: Version,
		Username: Username,
		PhoneNumber: PhoneNumber,
	}
}


// CreateUser создаёт нового пользователя с автоматически сгенерированными
// ID (UUID v4) и начальной версией 1.
func NewSaveUser(
	id 			uuid.UUID,
	username 	string,
	phoneNumber string,
) SaveUser {
	var (
		version = 1
	)

	return SaveUser{
		id,
		version,
		username,
		phoneNumber,
	}
}

func (u SaveUser) Validate() error {
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