package auth_service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error) {
	op := "Auth.Service.HashPassword"
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%v: %v", op, err)
	}
	return string(hashBytes), nil
}


func VerifyPassword(hashedPassword, providePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providePassword))
}