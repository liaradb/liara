package raw

import (
	"bufio"
	"bytes"
	"testing"
	"time"
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

func TestTime_WriteDataReadData(t *testing.T) {
	t.Parallel()

	o := NewTime(time.Now())

	data := make([]byte, TimeSize+2)
	data0 := o.WriteData(data)

	if l := len(data0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	o1 := Time{}
	data1 := o1.ReadData(data)
	if l := len(data1); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	if !o1.Equal(o) {
		t.Errorf("incorrect result: %v, expected: %v", o1, o)
	}
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
