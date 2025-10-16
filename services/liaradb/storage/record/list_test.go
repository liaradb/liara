package record

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"testing"
)

// TODO: Should this just test Page?
func TestList(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	l := List{}
	var position Offset = 20

	position -= 2
	if i := l.Add(position, 2); i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	position -= 4
	if i := l.Add(position, 4); i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 1)
	}

	position -= 6
	if i := l.Add(position, 6); i != 2 {
		t.Errorf("incorrect index: %v, expected: %v", i, 2)
	}

	if err := l.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := l.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var l2 List
	if err := l2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(l, l2) {
		t.Errorf("incorrect value: %v, expected: %v", l2, l)
	}
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
