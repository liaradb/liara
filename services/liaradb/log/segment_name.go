package log

import (
	"fmt"
	"regexp"
	"strconv"
)

var segmentRegexp = regexp.MustCompile("segment_([0-9]*)_([0-9]*).lr")

type SegmentName struct {
	id  SegmentID
	lsn LogSequenceNumber
}

func NewLogSegmentName(id SegmentID, lsn LogSequenceNumber) SegmentName {
	return SegmentName{
		id:  id,
		lsn: lsn,
	}
}

func ParseLogSegmentName(value string) SegmentName {
	matches := segmentRegexp.FindStringSubmatch(value)
	if len(matches) < 3 {
		return SegmentName{}
	}

	i, _ := strconv.Atoi(matches[1])
	l, _ := strconv.Atoi(matches[2])

	return NewLogSegmentName(SegmentID(i), LogSequenceNumber(l))
}

func (sn SegmentName) ID() SegmentID                        { return sn.id }
func (sn SegmentName) LogSequenceNumber() LogSequenceNumber { return sn.lsn }

func (sn SegmentName) String() string {
	return fmt.Sprintf("segment_%03v_%03v.lr", sn.id, sn.lsn)
}
