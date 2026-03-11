package serializer

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestWriteAllReadAll(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	a := base.String("a")
	b := base.String("b")

	if err := WriteAll(w, a, b); err != nil {
		t.Fatal(err)
	}

	var a2 base.String
	var b2 base.String

	if err := ReadAll(r, &a2, &b2); err != nil {
		t.Fatal(err)
	}

	if a2 != a {
		t.Errorf("incorrect result: %v, expected: %v", a2, a)
	}

	if b2 != b {
		t.Errorf("incorrect result: %v, expected: %v", b2, b)
	}
}

func TestWriteAllReadAll__Error(t *testing.T) {
	t.Parallel()

	buf := buffer.New(0)

	a := base.String("a")
	b := base.String("b")

	if err := WriteAll(buf, a, b); err != io.ErrShortWrite {
		t.Fatal("should return error")
	}

	var a2 base.String
	var b2 base.String

	if err := ReadAll(buf, &a2, &b2); err != io.EOF {
		t.Fatal("should return error")
	}
}

func TestSize(t *testing.T) {
	t.Parallel()

	a := base.String("a")
	b := base.String("b")

	want := a.Size() + b.Size()

	if s := Size(a, b); s != want {
		t.Errorf("incorrect size: %v, expected: %v", s, want)
	}
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
