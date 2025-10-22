package record

import "github.com/liaradb/liaradb/raw"

type TransactionID struct {
	baseUint64
}

const TransactionIDSize = raw.BaseUint64Size

func NewTransactionID(value uint64) TransactionID {
	return TransactionID{baseUint64(value)}
}
