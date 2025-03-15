package txmanager

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/madevara24/go-common/logger"
)

type contextDBType string

var ContextTxValue contextDBType = "TX"

// ExtractDB is used by other repo to extract the trx from context
func ExtractTx(ctx context.Context) (*sqlx.Tx, error) {

	db, ok := ctx.Value(ContextTxValue).(*sqlx.Tx)
	if !ok {
		return nil, ErrTxNotFound
	}

	return db, nil
}

func DBTransactionWrapperWithContext(ctx context.Context, db *sqlx.DB, closureFunc func(txCtx context.Context) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		logger.Log.Error(ctx, "failed when begin transaction", err)
		return ErrDbTransactionWrapper
	}
	txCtx := context.WithValue(ctx, ContextTxValue, tx)

	err = closureFunc(txCtx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Log.Error(ctx, "failed when rollback transaction", err)
			return ErrRollbackTransaction
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error(ctx, "failed when commit transaction", err)
		return ErrCommitTransaction
	}

	return nil
}
