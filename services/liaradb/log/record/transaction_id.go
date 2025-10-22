package record

type TransactionID struct {
	baseUint64
}

const TransactionIDSize = baseUint64Size

func NewTransactionID(value uint64) TransactionID {
	return TransactionID{baseUint64(value)}
}
