package log

import (
	"fmt"
	"regexp"
	"strconv"
)

var r = regexp.MustCompile("segment_([0-9]*)_([0-9]*).lr")

type LogSegmentName string

func NewLogSegmentName(index int, lsn LogSequenceNumber) LogSegmentName {
	return LogSegmentName(fmt.Sprintf("segment_%03v_%03v.lr", index, lsn))
}

func (lsn LogSegmentName) String() string { return string(lsn) }

func (lsn LogSegmentName) Value() (int, LogSequenceNumber) {
	matches := r.FindStringSubmatch(string(lsn))
	if len(matches) < 3 {
		return 0, 0
	}

	i, _ := strconv.Atoi(matches[1])
	l, _ := strconv.Atoi(matches[2])

	return i, LogSequenceNumber(l)
}
