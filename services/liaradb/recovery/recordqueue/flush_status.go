package recordqueue

import "github.com/liaradb/liaradb/recovery/logpage"

type flushStatus struct {
	highWater logpage.LogSequenceNumber
	lowWater  logpage.LogSequenceNumber
}

func (lf *flushStatus) HighWater() logpage.LogSequenceNumber { return lf.highWater }
func (lf *flushStatus) LowWater() logpage.LogSequenceNumber  { return lf.lowWater }
func (lf *flushStatus) isDirty() bool                        { return lf.lowWater != lf.highWater }

func (lf *flushStatus) init(lw, hw logpage.LogSequenceNumber) {
	lf.lowWater = lw
	lf.highWater = hw
}

func (lf *flushStatus) setHighWater(hw logpage.LogSequenceNumber) {
	lf.highWater = hw
}

func (lf *flushStatus) completeFlush() {
	lf.lowWater = lf.highWater
}
