package log

import (
	"fmt"
	"regexp"
	"strconv"
)

var logSegmentRegexp = regexp.MustCompile("segment_([0-9]*)_([0-9]*).lr")

type LogSegmentName struct {
	index             int
	logSequenceNumber LogSequenceNumber
}

func NewLogSegmentName(index int, lsn LogSequenceNumber) LogSegmentName {
	return LogSegmentName{
		index:             index,
		logSequenceNumber: lsn,
	}
}

func (lsn LogSegmentName) Index() int                           { return lsn.index }
func (lsn LogSegmentName) LogSequenceNumber() LogSequenceNumber { return lsn.logSequenceNumber }

func ParseLogSegmentName(value string) LogSegmentName {
	matches := logSegmentRegexp.FindStringSubmatch(value)
	if len(matches) < 3 {
		return LogSegmentName{}
	}

	i, _ := strconv.Atoi(matches[1])
	l, _ := strconv.Atoi(matches[2])

	return NewLogSegmentName(i, LogSequenceNumber(l))
}

func (lsn LogSegmentName) String() string {
	return fmt.Sprintf("segment_%03v_%03v.lr", lsn.index, lsn.logSequenceNumber)
}
