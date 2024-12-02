package service

import "context"

type TransactionRepository interface {
	Run(context.Context, func(tx Transaction) error) error
}

type Transaction interface {
	Commit() error
	Rollback() error
}
