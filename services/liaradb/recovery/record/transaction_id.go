package record

import "github.com/liaradb/liaradb/encoder/base"

type TransactionID struct {
	baseUint64
}

const TransactionIDSize = base.BaseUint64Size

func NewTransactionID(value uint64) TransactionID {
	return TransactionID{baseUint64(value)}
}
