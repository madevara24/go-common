package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(config *Configuration) error {
	sqldb, err := sql.Open("postgres", config.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	db := sqlx.NewDb(sqldb, "postgres")
	store.Store(config.SqlxKey, db)

	return nil
}

func DBTransactionWrapper(ctx context.Context, db *sqlx.DB, closureFunc func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = closureFunc(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
