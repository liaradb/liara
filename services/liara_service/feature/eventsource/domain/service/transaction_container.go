package service

import "context"

type TransactionContainer interface {
	Run(context.Context, func() error) error
}
