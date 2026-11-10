package service

import (
	"sync"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/transaction/record"
	"github.com/liaradb/liaradb/util/queue"
)

type tableVersion struct {
	mux     sync.RWMutex
	version value.GlobalVersion
	lw      value.GlobalVersion
	queue   queue.MapQueue[record.TransactionID, value.GlobalVersion]
}

func (tv *tableVersion) lowWater() value.GlobalVersion {
	tv.mux.RLock()
	defer tv.mux.RUnlock()

	return tv.lw
}

func (tv *tableVersion) setValue(v value.GlobalVersion) {
	tv.mux.Lock()
	defer tv.mux.Unlock()

	tv.version = v
}

func (tv *tableVersion) increment(txid record.TransactionID) value.GlobalVersion {
	tv.mux.Lock()
	defer tv.mux.Unlock()

	tv.version = tv.version.Increment()
	tv.queue.Push(txid, tv.version)
	return tv.version
}

func (lwv *tableVersion) complete(txid record.TransactionID) {
	lwv.mux.Lock()
	defer lwv.mux.Unlock()

	first := lwv.queue.IsFirst(txid)
	v, ok := lwv.queue.Remove(txid)
	if first && ok {
		lwv.lw = v
	}
}
