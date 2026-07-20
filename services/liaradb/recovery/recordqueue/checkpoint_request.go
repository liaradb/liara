package recordqueue

import (
	"context"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/async"
)

type CheckpointRequest = async.Request[CheckpointValue, logpage.LogSequenceNumber]

type CheckpointHandler struct {
	reqs async.Handler[CheckpointValue, logpage.LogSequenceNumber]
}

func NewCheckpointHandler() CheckpointHandler {
	return CheckpointHandler{
		reqs: make(async.Handler[CheckpointValue, logpage.LogSequenceNumber]),
	}
}

func (h CheckpointHandler) Reqs() async.Handler[CheckpointValue, logpage.LogSequenceNumber] {
	return h.reqs
}

func (h CheckpointHandler) Append(
	ctx context.Context,
	txids []record.TransactionID,
	time time.Time,
) (logpage.LogSequenceNumber, error) {
	return h.reqs.Send(ctx, CheckpointValue{
		txids: txids,
		time:  time,
	})
}

type CheckpointValue struct {
	txids []record.TransactionID
	time  time.Time
}

func (cv *CheckpointValue) Record(lsn logpage.LogSequenceNumber) *record.Record {
	return record.New(
		lsn,
		value.TenantID{},
		record.TransactionID{},
		record.NewTime(cv.time),
		record.ActionCheckpoint,
		record.CollectionSystem,
		cv.txIDsToData(),
		nil)
}

func (cv *CheckpointValue) txIDsToData() []byte {
	data := make([]byte, len(cv.txids)*record.TransactionIDSize)

	data0 := data
	for _, txid := range cv.txids {
		// There will always be enough space
		data0, _ = txid.WriteData(data0)
	}

	return data
}
