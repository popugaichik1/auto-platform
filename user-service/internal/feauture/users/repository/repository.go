package repository

import (
	core_logger "user-service/internal/core/logger"
	core_postgres_pool "user-service/internal/core/repository/postgres/pool"
)


type Repo struct {
	pool 	core_postgres_pool.Pool
	log 	*core_logger.Logger
}


func NewRepo(
	pool core_postgres_pool.Pool,
	log *core_logger.Logger,
) *Repo {
	return &Repo{
		pool: pool,
		log:  log,
	}
}

