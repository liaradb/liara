package value

import (
	"time"

	"github.com/liaradb/liaradb/encoder/raw"
)

type Time struct {
	baseTime
}

func NewTime(t time.Time) Time {
	return Time{raw.NewTime(t)}
}
