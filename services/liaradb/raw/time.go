package raw

import (
	"io"
	"time"
)

type Time struct {
	time.Time
}

const TimeSize = 8

func NewTime(t time.Time) Time {
	return Time{t}
}

func (Time) Size() int { return TimeSize }

func (t Time) Write(w io.Writer) error {
	return WriteInt64(w, t.Time.UnixMicro())
}

func (t *Time) Read(r io.Reader) error {
	var v int64
	if err := ReadInt64(r, &v); err != nil {
		return err
	}

	t.Time = time.UnixMicro(v)
	return nil
}

func (t Time) Equal(b Time) bool {
	return t.Time.Equal(b.Time)
}
