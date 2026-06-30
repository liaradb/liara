package recovery

import (
	"github.com/liaradb/liaradb/recovery/pageiterator"
	"github.com/liaradb/liaradb/recovery/record"
)

type logFlusher struct {
	highWater record.LogSequenceNumber
	lowWater  record.LogSequenceNumber
}

func (lf *logFlusher) HighWater() record.LogSequenceNumber { return lf.highWater }
func (lf *logFlusher) LowWater() record.LogSequenceNumber  { return lf.lowWater }
func (lf *logFlusher) isDirty() bool                       { return lf.lowWater != lf.highWater }

func (lf *logFlusher) initHighWater(it *pageiterator.PageIterator) error {
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

func (lf *logFlusher) setHighWater(hw record.LogSequenceNumber) {
	lf.highWater = hw
}

func (lf *logFlusher) completeFlush() {
	lf.lowWater = lf.highWater
}
