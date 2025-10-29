package transaction

import (
	"errors"
	"fmt"

	"github.com/liaradb/liaradb/recovery/record"
)

var (
	ErrTransactionFailed = errors.New("transaction failed")
)

func errTransactionFailed(txnum record.TransactionID, err error) error {
	return fmt.Errorf("transaction %v: %w: %w", txnum, ErrTransactionFailed, err)
}
