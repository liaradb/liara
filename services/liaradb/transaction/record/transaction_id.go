package record

import (
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/encoder/scan"
)

type TransactionID struct {
	baseUint64
}

const TransactionIDSize = base.Uint64Size

func NewTransactionID(value uint64) TransactionID {
	return TransactionID{baseUint64(value)}
}

// TODO: Should we use these?
func (b *TransactionID) ReadData(data []byte) []byte {
	v, data0 := scan.Uint64(data)
	*b = NewTransactionID(v)
	return data0
}

func (b TransactionID) WriteData(data []byte) []byte {
	return scan.SetUint64(data, b.Value())
}
