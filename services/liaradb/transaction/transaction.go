package transaction

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/log"
	"github.com/liaradb/liaradb/log/record"
)

type Transaction struct {
	id  record.TransactionID
	lsn record.LogSequenceNumber
	log *log.Log
}

func newTransaction(
	id record.TransactionID,
	log *log.Log,
) *Transaction {
	return &Transaction{
		id:  id,
		log: log,
	}
}

func (t Transaction) ID() record.TransactionID                     { return t.id }
func (t *Transaction) LogSequenceNumber() record.LogSequenceNumber { return t.lsn }

func (t *Transaction) Insert(ctx context.Context, now time.Time, data []byte) error {
	lsn, err := t.log.Append(ctx, t.id, now, record.ActionInsert, data, nil)
	if err != nil {
		return err
	}

	t.lsn = lsn
	return nil
}

func (t *Transaction) Commit(ctx context.Context, now time.Time) error {
	lsn, err := t.log.Append(ctx, t.id, now, record.ActionCommit, nil, nil)
	if err != nil {
		return err
	}

	t.lsn = lsn
	// TODO: Is this correct?
	return t.log.Flush(ctx, lsn)
}
