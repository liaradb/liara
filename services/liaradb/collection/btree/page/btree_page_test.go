package page

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/raw"
)

const (
	headerSize = 2 + btreePageHeaderSize
)

func TestBTreePage_Append(t *testing.T) {
	t.Parallel()

	const (
		size int16 = 256
		s0         = size - itemSize - headerSize
		s1         = s0 - itemSize - 16
		s2         = s1 - itemSize - 16
	)

	p := New(make([]byte, size))
	v0 := []byte{1, 2, 3, 4, 5}
	v1 := []byte{6, 7, 8, 9, 10}

	if s := p.Space(); s != s0 {
		t.Errorf("incorrect space: %v, expected: %v", s, s0)
	}

	i, b0, ok := p.Append(16)
	if !ok {
		t.Error("should get a buffer")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if s := p.Space(); s != s1 {
		t.Errorf("incorrect space: %v, expected: %v", s, s1)
	}

	if _, err := raw.NewBufferFromSlice(b0).Write(v0); err != nil {
		t.Error(err)
	}

	i, b1, ok := p.Append(16)
	if !ok {
		t.Error("should get a buffer")
	} else if i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 1)
	}

	if s := p.Space(); s != s2 {
		t.Errorf("incorrect space: %v, expected: %v", s, s2)
	}

	if _, err := raw.NewBufferFromSlice(b1).Write(v1); err != nil {
		t.Error(err)
	}

	// if _, err := raw.NewBufferFromSlice(b0).Seek(0, io.SeekStart); err != nil {
	// 	t.Error(err)
	// }

	r0 := make([]byte, 5)
	if _, err := raw.NewBufferFromSlice(b0).Read(r0); err != nil {
		t.Error(err)
	}

	if !slices.Equal(r0, v0) {
		t.Errorf("incorrect result: %v, expected: %v", r0, v0)
	}

	// if _, err := raw.NewBufferFromSlice(b1).Seek(0, io.SeekStart); err != nil {
	// 	t.Error(err)
	// }

	r1 := make([]byte, 5)
	if _, err := raw.NewBufferFromSlice(b1).Read(r1); err != nil {
		t.Error(err)
	}

	if !slices.Equal(r1, v1) {
		t.Errorf("incorrect result: %v, expected: %v", r1, v1)
	}
}

func TestBTreePage_Insert(t *testing.T) {
	t.Parallel()

	const (
		size int16 = 256
		s0         = size - itemSize - headerSize
		s1         = s0 - itemSize - 16
		s2         = s1 - itemSize - 16
	)

	p := New(make([]byte, size))
	v0 := []byte{1, 2, 3, 4, 5}
	v1 := []byte{6, 7, 8, 9, 10}

	if s := p.Space(); s != s0 {
		t.Fatalf("incorrect space: %v, expected: %v", s, s0)
	}

	i, b0, ok := p.Insert(16, 0)
	if !ok {
		t.Fatal("should get a buffer")
	} else if i != 0 {
		t.Fatalf("incorrect index: %v, expected: %v", i, 0)
	}

	if s := p.Space(); s != s1 {
		t.Fatalf("incorrect space: %v, expected: %v", s, s1)
	}

	if _, err := raw.NewBufferFromSlice(b0).Write(v0); err != nil {
		t.Fatal(err)
	}

	i, b1, ok := p.Insert(16, 1)
	if !ok {
		t.Fatal("should get a buffer")
	} else if i != 1 {
		t.Fatalf("incorrect index: %v, expected: %v", i, 1)
	}

	if s := p.Space(); s != s2 {
		t.Fatalf("incorrect space: %v, expected: %v", s, s2)
	}

	if _, err := raw.NewBufferFromSlice(b1).Write(v1); err != nil {
		t.Fatal(err)
	}

	// if _, err := raw.NewBufferFromSlice(b0).Seek(0, io.SeekStart); err != nil {
	// 	t.Fatal(err)
	// }

	r0 := make([]byte, 5)
	if _, err := raw.NewBufferFromSlice(b0).Read(r0); err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(r0, v0) {
		t.Fatalf("incorrect result: %v, expected: %v", r0, v0)
	}

	// if _, err := raw.NewBufferFromSlice(b1).Seek(0, io.SeekStart); err != nil {
	// 	t.Fatal(err)
	// }

	r1 := make([]byte, 5)
	if _, err := raw.NewBufferFromSlice(b1).Read(r1); err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(r1, v1) {
		t.Errorf("incorrect result: %v, expected: %v", r1, v1)
	}
}

func TestBTreePage_Space(t *testing.T) {
	t.Parallel()

	p := New(make([]byte, 16+itemSize+headerSize))

	if s := p.Space(); s != 16 {
		t.Errorf("incorrect space: %v, expected: %v", s, 16)
	}

	if _, _, ok := p.Append(16); !ok {
		t.Error("should get a buffer")
	}

	if s := p.Space(); s != 0 {
		t.Errorf("incorrect space: %v, expected: %v", s, 0)
	}

	if _, _, ok := p.Append(16); ok {
		t.Error("should not get a buffer")
	}
}

func TestBTreePage_Child(t *testing.T) {
	t.Parallel()

	p := New(make([]byte, 256))
	values := [][]byte{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10}}
	_, b0, ok := p.Append(16)
	if !ok {
		t.Error("should get a buffer")
	}

	if _, err := raw.NewBufferFromSlice(b0).Write(values[0]); err != nil {
		t.Error(err)
	}

	_, b1, ok := p.Append(16)
	if !ok {
		t.Error("should get a buffer")
	}

	if _, err := raw.NewBufferFromSlice(b1).Write(values[1]); err != nil {
		t.Error(err)
	}

	result := make([][]byte, 0, 2)
	for i := range 2 {
		c, ok := p.Child(int16(i))
		if !ok {
			t.Fatal("should get a buffer")
		}

		v := make([]byte, 5)
		if _, err := raw.NewBufferFromSlice(c).Read(v); err != nil {
			t.Fatal(err)
		}

		result = append(result, v)
	}

	if !slices.EqualFunc(result, values, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, values)
	}
}

func TestBTreePage_Children(t *testing.T) {
	t.Parallel()

	p := New(make([]byte, 256))
	values := [][]byte{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10}}
	_, b0, ok := p.Append(16)
	if !ok {
		t.Error("should get a buffer")
	}

	if _, err := raw.NewBufferFromSlice(b0).Write(values[0]); err != nil {
		t.Error(err)
	}

	_, b1, ok := p.Append(16)
	if !ok {
		t.Error("should get a buffer")
	}

	if _, err := raw.NewBufferFromSlice(b1).Write(values[1]); err != nil {
		t.Error(err)
	}

	result := make([][]byte, 0, 2)
	for c := range p.Children() {
		v := make([]byte, 5)
		if _, err := raw.NewBufferFromSlice(c).Read(v); err != nil {
			t.Fatal(err)
		}

		result = append(result, v)
	}

	if !slices.EqualFunc(result, values, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, values)
	}
}
