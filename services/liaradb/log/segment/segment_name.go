package segment

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/liaradb/liaradb/log/record"
)

var segmentRegexp = regexp.MustCompile("segment_([0-9a-f]*)_([0-9a-f]*).lr")

type SegmentName struct {
	id  SegmentID
	lsn record.LogSequenceNumber
}

func NewSegmentName(id SegmentID, lsn record.LogSequenceNumber) SegmentName {
	return SegmentName{
		id:  id,
		lsn: lsn,
	}
}

func ParseSegmentName(value string) SegmentName {
	matches := segmentRegexp.FindStringSubmatch(value)
	if len(matches) < 3 {
		return SegmentName{}
	}

	i, _ := strconv.ParseUint(matches[1], 16, 64)
	l, _ := strconv.ParseUint(matches[2], 16, 64)

	return NewSegmentName(SegmentID(i), record.NewLogSequenceNumber(l))
}

func (sn SegmentName) ID() SegmentID                               { return sn.id }
func (sn SegmentName) LogSequenceNumber() record.LogSequenceNumber { return sn.lsn }

func (sn SegmentName) Next(lsn record.LogSequenceNumber) SegmentName {
	return SegmentName{
		id:  sn.id + 1,
		lsn: lsn,
	}
}

func (sn SegmentName) String() string {
	return fmt.Sprintf("segment_%016x_%016x.lr", sn.id, sn.lsn.Value())
}
