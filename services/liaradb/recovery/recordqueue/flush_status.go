package recordqueue

import "github.com/liaradb/liaradb/recovery/record"

type flushStatus struct {
	highWater record.LogSequenceNumber
	lowWater  record.LogSequenceNumber
}

func (lf *flushStatus) HighWater() record.LogSequenceNumber { return lf.highWater }
func (lf *flushStatus) LowWater() record.LogSequenceNumber  { return lf.lowWater }
func (lf *flushStatus) isDirty() bool                       { return lf.lowWater != lf.highWater }

func (lf *flushStatus) init(lw, hw record.LogSequenceNumber) {
	lf.lowWater = lw
	lf.highWater = hw
}

func (lf *flushStatus) setHighWater(hw record.LogSequenceNumber) {
	lf.highWater = hw
}

func (lf *flushStatus) completeFlush() {
	lf.lowWater = lf.highWater
}
