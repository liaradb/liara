package link

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/liaradb/liaradb/encoder/page"
)

func TestFilePosition(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var p FilePosition = 123456
	if err := p.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := p.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var p2 FilePosition
	if err := p2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if p != p2 {
		t.Errorf("incorrect value: %v, expected: %v", p2, p)
	}
}

func TestFilePosition_Offset(t *testing.T) {
	var p FilePosition = 123
	want := page.Offset(2 * 123)
	if o := p.Offset(2); o != want {
		t.Errorf("incorrect offset: %v, expected: %v", o, want)
	}
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
