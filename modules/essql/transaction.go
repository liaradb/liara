package essql

import (
	"context"
	"database/sql"
	"errors"
)

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
