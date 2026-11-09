package recordqueue

import (
	"sync"

	"github.com/liaradb/liaradb/recovery/logpage"
)

type flushStatus struct {
	highWater logpage.LogSequenceNumber
	lowWater  logpage.LogSequenceNumber
	hwMux     sync.RWMutex
	lwMux     sync.RWMutex
}

func (lf *flushStatus) HighWater() logpage.LogSequenceNumber {
	lf.hwMux.RLock()
	defer lf.hwMux.RUnlock()

	return lf.highWater
}

func (lf *flushStatus) LowWater() logpage.LogSequenceNumber {
	lf.lwMux.RLock()
	defer lf.lwMux.RUnlock()

	return lf.lowWater
}

func (lf *flushStatus) isDirty() bool {
	lf.lwMux.RLock()
	lf.hwMux.RLock()
	defer lf.lwMux.RUnlock()
	defer lf.hwMux.RUnlock()

	return lf.lowWater != lf.highWater
}

func (lf *flushStatus) init(lw, hw logpage.LogSequenceNumber) {
	lf.lwMux.Lock()
	lf.hwMux.Lock()
	defer lf.lwMux.Unlock()
	defer lf.hwMux.Unlock()

	lf.lowWater = lw
	lf.highWater = hw
}

func (lf *flushStatus) setHighWater(hw logpage.LogSequenceNumber) {
	lf.hwMux.Lock()
	defer lf.hwMux.Unlock()

	lf.highWater = hw
}

func (lf *flushStatus) completeFlush() {
	lf.lwMux.Lock()
	defer lf.lwMux.Unlock()

	lf.lowWater = lf.highWater
}
