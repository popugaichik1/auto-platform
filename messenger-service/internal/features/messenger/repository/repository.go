package repository

import (
	core_postgres_pool "messenger-service/internal/core/repository/postgres/pool"
)

type Repository struct {
	pool core_postgres_pool.Pool
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{pool: pool}
}
