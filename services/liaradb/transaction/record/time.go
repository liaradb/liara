package record

import (
	"time"

	"github.com/liaradb/liaradb/encoder/base"
)

type Time struct {
	baseTime
}

func NewTime(t time.Time) Time {
	return Time{base.NewTime(t)}
}

func (t Time) Equal(b Time) bool {
	return t.baseTime.Equal(b.baseTime)
}
