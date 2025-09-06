package log

type Log struct {
	highWater LogSequenceNumber
	lowWater  LogSequenceNumber
}

type LogSequenceNumber uint64

func (l *Log) HighWater() LogSequenceNumber { return l.highWater }
func (l *Log) LowWater() LogSequenceNumber  { return l.lowWater }

func (l *Log) Append(data []byte) (LogSequenceNumber, error) {
	l.highWater++
	return l.highWater, nil
}

func (l *Log) Flush(lsn LogSequenceNumber) error {
	lsn = min(lsn, l.highWater)
	l.lowWater = lsn
	return nil
}
