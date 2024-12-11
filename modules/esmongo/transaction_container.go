package esmongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionContainer struct {
	db *mongo.Database
}

func NewTransactionContainer(
	db *mongo.Database,
) *TransactionContainer {
	return &TransactionContainer{
		db: db,
	}
}

func (tc *TransactionContainer) Transaction(ctx context.Context, tx func() error) (err error) {
	session, err := tc.db.Client().StartSession()
	if err != nil {
		return err
	}

	if err := session.StartTransaction(); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, session.AbortTransaction(ctx))
		} else {
			err = session.CommitTransaction(ctx)
		}
	}()

	err = tx()

	return
}
