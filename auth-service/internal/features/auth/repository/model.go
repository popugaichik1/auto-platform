package auth_postgres_repository

import (
	"time"

	"github.com/google/uuid"
	core_domain "github.com/zosinkin/social_network/internal/core/domain"
	core_postgres_pool "github.com/zosinkin/social_network/internal/core/repository/postgres/pool"
)

type UserModel struct {
	ID 				uuid.UUID
	Username 		string
	PhoneNumber 	string
	PasswordHash    string
}


type RefreshTokenModel struct {
	ID  			uuid.UUID
	UserID 			uuid.UUID
	Token 			string
	ExpiresAt 		time.Time
	CreatedAt 		time.Time
	Revoked 		bool
}


// Scan заполняет поля модели из результата запроса к БД.
func (m *UserModel) Scan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Username,
		&m.PhoneNumber,
	)
}


func (m *UserModel) ScanWithPassword(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.Username,
		&m.PhoneNumber,
		&m.PasswordHash,
	)
}

func (m *RefreshTokenModel) TokenScan(row core_postgres_pool.Row) error {
	return row.Scan(
		&m.ID,
		&m.UserID,
		&m.Token,
		&m.ExpiresAt,
		&m.CreatedAt,
		&m.Revoked,
	)
}


func modelToUserDomain(model UserModel) core_domain.AuthUser {
	return core_domain.AuthUser{
		ID:           model.ID,
		Username:     model.Username,
		PhoneNumber:  model.PhoneNumber,
		PasswordHash: model.PasswordHash,
	}
}

func modelTokenDomain(model RefreshTokenModel) core_domain.RefreshToken {
	return core_domain.RefreshToken{
		ID:        model.ID,
		UserID:    model.UserID,
		Token:     model.Token,
		ExpiresAt: model.ExpiresAt,
		CreatedAt: model.CreatedAt,
		Revoked:   model.Revoked,
	}
}