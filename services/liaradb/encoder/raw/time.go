package raw

import (
	"io"
	"time"

	"github.com/liaradb/liaradb/encoder/scan"
)

type baseTime = time.Time

type Time struct {
	baseTime
}

const TimeSize = 16

func NewTime(t time.Time) Time {
	return Time{t.UTC()}
}

func (t Time) Value() time.Time { return t.baseTime }
func (Time) Size() int          { return TimeSize }

func (t Time) Write(w io.Writer) error {
	if err := WriteInt64(w, t.baseTime.Unix()); err != nil {
		return err
	}

	return WriteInt64(w, int64(t.baseTime.Nanosecond()))
}

func (t *Time) Read(r io.Reader) error {
	var s int64
	if err := ReadInt64(r, &s); err != nil {
		return err
	}

	var n int64
	if err := ReadInt64(r, &n); err != nil {
		return err
	}

	t.baseTime = time.Unix(s, n).UTC()
	return nil
}

func (t *Time) WriteData(data []byte) []byte {
	data0 := scan.SetInt64(data, t.baseTime.Unix())
	return scan.SetInt64(data0, int64(t.baseTime.Nanosecond()))
}

func (t *Time) ReadData(data []byte) []byte {
	s, data0 := scan.Int64(data)
	n, data1 := scan.Int64(data0)
	t.baseTime = time.Unix(s, n).UTC()
	return data1
}

func (t Time) Equal(b Time) bool {
	return t.baseTime.Equal(b.baseTime)
}

func (t *Time) ReadWrite() {
	s := t.Unix()
	n := int64(t.Nanosecond())

	time.Unix(s, n)
}
