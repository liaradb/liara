package raw

import (
	"io"
	"time"
)

type baseTime = time.Time

type Time struct {
	baseTime
}

const TimeSize = 8

func NewTime(t time.Time) Time {
	return Time{t}
}

func (Time) Size() int { return TimeSize }

func (t Time) Write(w io.Writer) error {
	return WriteInt64(w, t.baseTime.UnixMicro())
}

func (t *Time) Read(r io.Reader) error {
	var v int64
	if err := ReadInt64(r, &v); err != nil {
		return err
	}

	t.baseTime = time.UnixMicro(v).UTC()
	return nil
}

func (t Time) Equal(b Time) bool {
	return t.baseTime.Equal(b.baseTime)
}
