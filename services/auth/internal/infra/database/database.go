package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(p *pgxpool.Pool) *Database {
	return &Database{pool: p}
}

func (d *Database) Execute(ctx context.Context, query string, args ...any) error {
	exc := d.executor(ctx)
	res, err := exc.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if ra := res.RowsAffected(); ra == 0 {
		return ce.ErrDBAffectNoRows
	}
	return nil
}

func (d *Database) Query(ctx context.Context, query string, args ...any) pgx.Row {
	exc := d.executor(ctx)
	return exc.QueryRow(ctx, query, args...)
}

func (d *Database) WithinTx(ctx context.Context) bool {
	return txFromCtx(ctx) != nil
}
