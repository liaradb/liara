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

	if s := p.String(); s != "123456" {
		t.Errorf("incorrect string: %v, expected: %v", s, "123456")
	}
}

func TestFilePosition_ReadDataWriteData(t *testing.T) {
	t.Parallel()

	fp := FilePosition(1)

	data := make([]byte, 12)
	data0, ok := fp.WriteData(data)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(data0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	fp0 := FilePosition(0)
	data1, ok := fp0.ReadData(data)
	if !ok {
		t.Error("unable to read")
	}

	if l := len(data1); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if fp0 != fp {
		t.Errorf("incorrect value: %v, expected: %v", fp0, fp)
	}
}

func TestFilePosition_Offset(t *testing.T) {
	t.Parallel()

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
