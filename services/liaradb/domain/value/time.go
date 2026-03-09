package value

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
