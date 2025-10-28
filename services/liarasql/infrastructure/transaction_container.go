package infrastructure

import (
	"context"
	"database/sql"
	"errors"
)

type TransactionContainer struct {
	db   *sql.DB
	opts *sql.TxOptions
}

func NewTransactionContainer(db *sql.DB, opts *sql.TxOptions) *TransactionContainer {
	return &TransactionContainer{
		db:   db,
		opts: opts,
	}
}

func (tr *TransactionContainer) Run(ctx context.Context, transaction func() error) error {
	return runTx(ctx, tr.db, tr.opts, transaction)
}

func runTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, transaction func() error) (err error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	err = transaction()

	return
}
