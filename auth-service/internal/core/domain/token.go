package core_domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID  			uuid.UUID
	UserID 			uuid.UUID
	Token 			string
	ExpiresAt 		time.Time
	CreatedAt 		time.Time
	Revoked 		bool
}


func NewRefreshToken(
	userID 		uuid.UUID,
	ttl  		time.Duration,
) RefreshToken {
	ID := uuid.New()
	expiresAt := time.Now().Add(ttl)
	tokenID := ID.String()
	createdAt := time.Now()
	return RefreshToken{
		ID: 		ID,
		UserID: 	userID,
		Token:	 	tokenID,
		ExpiresAt: 	expiresAt,
		CreatedAt: 	createdAt,
		Revoked: 	false,
	}

}

