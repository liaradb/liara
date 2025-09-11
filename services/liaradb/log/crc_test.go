package log

import (
	"io"
	"testing"
)

func TestCRC(t *testing.T) {
	r, w := createReaderWriter()

	var c CRC = NewCRC([]byte{1, 2, 3, 4, 5})
	if err := c.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var c2 CRC
	if err := c2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if c != c2 {
		t.Errorf("incorrect value: %v, expected: %v", c2, c)
	}
}

func TestCRC_Compare(t *testing.T) {
	a := []byte{1, 2, 3, 4, 5}
	b := []byte{6, 7, 8, 9, 0}

	c := NewCRC(a)

	if !c.Compare(a) {
		t.Error("should be true")
	}

	if c.Compare(b) {
		t.Error("should be false")
	}
}
