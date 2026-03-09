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

	tm2 := Time{}

	if err := tm2.Read(r); err != nil {
		t.Fatal(err)
	}

	if tm2 != tm {
		t.Errorf("incorrect result: %v, expected: %v", tm2, tm)
	}
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
