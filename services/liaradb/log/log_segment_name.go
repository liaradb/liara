package log

import "fmt"

type LogSegmentName string

func NewLogSegmentName(index int) LogSegmentName {
	return LogSegmentName(fmt.Sprintf("segment_%03v.lr", index))
}

func (lsn LogSegmentName) String() string { return string(lsn) }
