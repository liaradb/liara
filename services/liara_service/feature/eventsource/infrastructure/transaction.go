package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/service"
)

type TransactionRepository struct {
	db   *sql.DB
	opts *sql.TxOptions
}

func NewTransactionRepository(db *sql.DB, opts *sql.TxOptions) *TransactionRepository {
	return &TransactionRepository{
		db:   db,
		opts: opts,
	}
}

func (tr *TransactionRepository) Run(ctx context.Context, transaction func(tx service.Transaction) error) error {
	return runTx(ctx, tr.db, tr.opts, func(tx *sql.Tx) error {
		return transaction(tx)
	})
}

func runTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, transaction func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	err = transaction(tx)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}
