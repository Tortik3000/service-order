package transactor

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Tortik3000/service-order/pkg/postgres"
)

type Transactor interface {
	WithTx(ctx context.Context, function func(ctx context.Context) error) error
	GetConn(ctx context.Context) (postgres.Conn, error)
}

type transactor struct {
	pool *pgxpool.Pool
}

var _ Transactor = (*transactor)(nil)

func New(pool *pgxpool.Pool) *transactor {
	return &transactor{
		pool: pool,
	}
}

type (
	txKey         struct{}
	txRequiredKey struct{}
)

func (t *transactor) WithTx(
	ctx context.Context,
	function func(ctx context.Context) error,
) (txErr error) {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if txErr != nil {
			_ = tx.Rollback(ctx)
			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			txErr = err
		}
	}()

	ctxWithTx := context.WithValue(ctx, txKey{}, tx)
	ctxWithTx = context.WithValue(ctxWithTx, txRequiredKey{}, true)

	return function(ctxWithTx)
}

func (t *transactor) GetConn(
	ctx context.Context,
) (postgres.Conn, error) {
	required, err := t.isTransactionRequired(ctx)
	if err != nil {
		return nil, err
	}

	if required {
		tx := t.getTx(ctx)
		return tx, nil
	}

	return t.pool, nil
}

func (t *transactor) isTransactionRequired(ctx context.Context) (bool, error) {
	required, ok := ctx.Value(txRequiredKey{}).(bool)
	if ok {
		if required {
			tx := t.getTx(ctx)
			if tx == nil {
				return false, errors.New("transaction required, but not found")
			}
		}
		return required, nil
	}

	return false, nil
}

func (t *transactor) getTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}
