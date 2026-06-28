package core_pgx_pool

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_pool "messenger-service/internal/core/repository/postgres/pool"
)

type pgxRows struct{ pgx.Rows }

type pgxRow struct{ pgx.Row }

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}

	return nil
}

type pgxCommandTag struct{ pgconn.CommandTag }

type pgxTx struct{ pgx.Tx }

func (t pgxTx) Query(ctx context.Context, sql string, args ...any) (core_postgres_pool.Rows, error) {
	rows, err := t.Tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxRows{rows}, nil
}

func (t pgxTx) QueryRow(ctx context.Context, sql string, args ...any) core_postgres_pool.Row {
	return pgxRow{t.Tx.QueryRow(ctx, sql, args...)}
}

func (t pgxTx) Exec(ctx context.Context, sql string, args ...any) (core_postgres_pool.CommandTag, error) {
	tag, err := t.Tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxCommandTag{tag}, nil
}

func (t pgxTx) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

func (t pgxTx) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}

func mapErrors(err error) error {
	const pgxViolatesForeignKeyErrorCode = "23503"

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgxViolatesForeignKeyErrorCode {
			return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrViolatesForeignKey)
		}
	}

	return fmt.Errorf("%v: %w", err, core_postgres_pool.ErrUnknown)
}
