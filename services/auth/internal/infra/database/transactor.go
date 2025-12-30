package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/erteku/services/auth/internal/utils/ce"
)

type Transactor struct {
	pool *pgxpool.Pool
}

func NewTransactor(p *pgxpool.Pool) *Transactor {
	return &Transactor{pool: p}
}

func (t *Transactor) WithTx(ctx context.Context, fn func(context.Context) *ce.Error) *ce.Error {
	tx := txFromCtx(ctx)
	newTx := false

	var err error
	if tx == nil {
		tx, err = t.pool.Begin(ctx)
		if err != nil {
			return ce.NewError(ce.CodeDBTX, ce.MsgInternalServer, err)
		}

		ctx = txToCtx(ctx, tx)
		newTx = true
	}
	if err := fn(ctx); err != nil {
		if newTx {
			_ = tx.Rollback(ctx)
		}
		return err
	}
	if newTx {
		if err := tx.Commit(ctx); err != nil {
			return ce.NewError(ce.CodeDBTX, ce.MsgInternalServer, err)
		}
	}
	return nil
}
