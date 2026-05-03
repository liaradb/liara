package multi

import (
	"io"
	"testing"

	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestReaderWriter(t *testing.T) {
	t.Parallel()

	a, b, c := buffer.New(8), buffer.New(8), buffer.New(8)

	w := NewWriter(a, b, c)
	r := NewReader(a, b, c)

	want := make([]byte, 0, 24)
	for i := range cap(want) {
		want = append(want, byte(i))
	}

	if _, err := w.Write(want); err != nil {
		t.Fatal(err)
	}

	if _, err := a.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	if _, err := b.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	result := make([]byte, 24)
	if _, err := r.Read(result); err != nil {
		t.Fatal(err)
	}
}

func TestReaderWriter_Append(t *testing.T) {
	t.Parallel()

	a, b, c := buffer.New(8), buffer.New(8), buffer.New(8)

	w := NewWriter()
	w.Append(a)
	w.Append(b)
	w.Append(c)

	r := NewReader()
	r.Append(a)
	r.Append(b)
	r.Append(c)

	want := make([]byte, 0, 24)
	for i := range cap(want) {
		want = append(want, byte(i))
	}

	if _, err := w.Write(want); err != nil {
		t.Fatal(err)
	}

	if _, err := a.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	if _, err := b.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	result := make([]byte, 24)
	if _, err := r.Read(result); err != nil {
		t.Fatal(err)
	}
}

// TODO: Test this
func TestReaderWriter__Error(t *testing.T) {
	t.Parallel()
	t.Skip()

	a, b, c := buffer.New(7), buffer.New(8), buffer.New(8)

	w := NewWriter(a, b, c)
	r := NewReader(a, b, c)

	want := make([]byte, 0, 24)
	for i := range cap(want) {
		want = append(want, byte(i))
	}

	if _, err := w.Write(want); err != nil {
		t.Fatal(err)
	}

	if _, err := a.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	if _, err := b.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}

	result := make([]byte, 24)
	if _, err := r.Read(result); err != nil {
		t.Fatal(err)
	}
}
