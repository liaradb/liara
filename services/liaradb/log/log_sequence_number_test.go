package log

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestLogSequenceNumber(t *testing.T) {
	r, w := createReaderWriter()

	var lsn LogSequenceNumber = 123456
	if err := lsn.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var lsn2 LogSequenceNumber
	if err := lsn2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if lsn != lsn2 {
		t.Errorf("incorrect value: %v, expected: %v", lsn2, lsn)
	}
}

func createReaderWriter() (*bufio.Reader, *bufio.Writer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer),
		bufio.NewWriter(buffer)
}
