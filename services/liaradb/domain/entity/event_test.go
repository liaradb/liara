package entity

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
)

func TestEvent(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var e = Event{
		GlobalVersion: value.NewGlobalVersion(1),
		Data:          value.NewData([]byte{}),
	}
	if err := e.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len() - raw.HeaderSize
	if s := e.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var e2 Event
	if err := e2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	// Data comparison doesn't allow nil slice
	if !reflect.DeepEqual(e, e2) {
		t.Errorf("incorrect value: %v, expected: %v", e2, e)
	}
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
