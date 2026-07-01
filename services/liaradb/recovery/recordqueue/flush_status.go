package recordqueue

import (
	"github.com/liaradb/liaradb/recovery/pageiterator"
	"github.com/liaradb/liaradb/recovery/record"
)

type flushStatus struct {
	highWater record.LogSequenceNumber
	lowWater  record.LogSequenceNumber
}

func (lf *flushStatus) HighWater() record.LogSequenceNumber { return lf.highWater }
func (lf *flushStatus) LowWater() record.LogSequenceNumber  { return lf.lowWater }
func (lf *flushStatus) isDirty() bool                       { return lf.lowWater != lf.highWater }

func (lf *flushStatus) initHighWater(it *pageiterator.PageIterator) error {
	lf.lowWater = record.NewLogSequenceNumber(0)
	lf.highWater = record.NewLogSequenceNumber(0)

	hw := false
	for rc, err := range it.Reverse() {
		if err != nil {
			return err
		}

		if !hw {
			lf.highWater = rc.LogSequenceNumber()
			hw = true
		}

		if rc.Action() == record.ActionCheckpoint {
			lf.lowWater = rc.LogSequenceNumber()
			break
		}
	}

	return nil
}

func (lf *flushStatus) setHighWater(hw record.LogSequenceNumber) {
	lf.highWater = hw
}

func (lf *flushStatus) completeFlush() {
	lf.lowWater = lf.highWater
}
