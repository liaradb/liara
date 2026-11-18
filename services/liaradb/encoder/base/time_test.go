package base

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestTime(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tm := NewTime(now)

	value := now.UTC()
	if v := tm.Value(); v != value {
		t.Errorf("incorrect value: %v, expected: %v", v, value)
	}
}

func TestTime_WriteRead(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tm := NewTime(now)

	r, w := newReaderWriter()
	if err := tm.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := tm.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	tm2 := Time{}

	if err := tm2.Read(r); err != nil {
		t.Fatal(err)
	}

	if tm2 != tm {
		t.Errorf("incorrect result: %v, expected: %v", tm2, tm)
	}
}

func TestTime_WriteRead__Error(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tm := NewTime(now)

	// Test first error
	buf := buffer.New(7)

	if err := tm.Write(buf); err == nil {
		t.Fatal("should return error")
	}

	buf.SeekStart()

	if err := tm.Read(buf); err == nil {
		t.Fatal("should return error")
	}

	// Test second error
	buf = buffer.New(15)

	if err := tm.Write(buf); err == nil {
		t.Fatal("should return error")
	}

	buf.SeekStart()

	if err := tm.Read(buf); err == nil {
		t.Fatal("should return error")
	}
}

func TestTime_Equal(t *testing.T) {
	now := time.Now()

	a := NewTime(now)
	b := NewTime(now)
	c := NewTime(now.Add(time.Second))

	if !a.Equal(b) {
		t.Error("should equal")
	}

	if a.Equal(c) {
		t.Error("should not equal")
	}
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
