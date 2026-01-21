package node

import (
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/storage/storagetesting"
)

const (
	testHeaderSize = 2 + headerSize
)

func TestNode_Append(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Append)
}

func testNode_Append(t *testing.T) {
	const (
		size int16 = 256
		s0         = size - itemSize - testHeaderSize
		s1         = s0 - itemSize - 16
		s2         = s1 - itemSize - 16
	)

	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	p := NewFromSlice(b.Raw())
	v0 := []byte{1, 2, 3, 4, 5}
	v1 := []byte{6, 7, 8, 9, 10}

	if s := p.Space(); s != s0 {
		t.Errorf("incorrect space: %v, expected: %v", s, s0)
	}

	i, ok := p.Append(v0)
	if !ok {
		t.Fatal("should append record")
	} else if i != 0 {
		t.Fatalf("incorrect index: %v, expected: %v", i, 0)
	}

	// if s := p.Space(); s != s1 {
	// 	t.Errorf("incorrect space: %v, expected: %v", s, s1)
	// }

	i, ok = p.Append(v1)
	if !ok {
		t.Fatal("should append record")
	} else if i != 1 {
		t.Fatalf("incorrect index: %v, expected: %v", i, 1)
	}

	// if s := p.Space(); s != s2 {
	// 	t.Errorf("incorrect space: %v, expected: %v", s, s2)
	// }

	// if _, err := raw.NewBufferFromSlice(b0).Seek(0, io.SeekStart); err != nil {
	// 	t.Error(err)
	// }

	result := make([][]byte, 0)
	for i := range p.Items() {
		result = append(result, i)
	}
	// // r0 := make([]byte, 5)
	// // if _, err := raw.NewBufferFromSlice(b0).Read(r0); err != nil {
	// // 	t.Error(err)
	// // }

	r0 := result[0]
	if !slices.Equal(r0, v0) {
		t.Errorf("incorrect result: %v, expected: %v", r0, v0)
	}

	// if _, err := raw.NewBufferFromSlice(b1).Seek(0, io.SeekStart); err != nil {
	// 	t.Error(err)
	// }

	// r1 := make([]byte, 5)
	// if _, err := raw.NewBufferFromSlice(b1).Read(r1); err != nil {
	// 	t.Error(err)
	// }

	r1 := result[1]
	if !slices.Equal(r1, v1) {
		t.Errorf("incorrect result: %v, expected: %v", r1, v1)
	}
}

func TestNode_Space(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Space)
}

func testNode_Space(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 16+itemSize+testHeaderSize)
	b := createBuffer(t, s)
	defer b.Release()

	p := NewFromSlice(b.Raw())

	if s := p.Space(); s != 16 {
		t.Fatalf("incorrect space: %v, expected: %v", s, 16)
	}

	if _, ok := p.Append(make([]byte, 16)); !ok {
		t.Fatal("should append record")
	}

	if s := p.Space(); s != 0 {
		t.Fatalf("incorrect space: %v, expected: %v", s, 0)
	}

	if _, ok := p.Append(make([]byte, 16)); ok {
		t.Fatal("should not get a buffer")
	}
}

func TestNode_Child(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Child)
}

func testNode_Child(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	b.Release()

	p := NewFromSlice(b.Raw())
	values := [][]byte{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10}}

	if _, ok := p.Append(values[0]); !ok {
		t.Fatal("should append record")
	}

	if _, ok := p.Append(values[1]); !ok {
		t.Fatal("should append record")
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

func TestNode_Items(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Items)
}

func testNode_Items(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	p := NewFromSlice(b.Raw())
	values := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10}}

	for _, v := range values {
		if _, ok := p.Append(v); !ok {
			t.Fatal("should append record")
		}
	}

	result := make([][]byte, 0, len(values))
	for c := range p.Items() {
		v := make([]byte, 2)
		if _, err := raw.NewBufferFromSlice(c).Read(v); err != nil {
			t.Fatal(err)
		}

		result = append(result, v)
	}

	if !slices.EqualFunc(result, values, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, values)
	}
}

func TestNode_ChildrenRange(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_ChildrenRange)
}

func testNode_ChildrenRange(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	p := NewFromSlice(b.Raw())
	values := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10}}

	for _, v := range values {
		if _, ok := p.Append(v); !ok {
			t.Fatal("should append record")
		}
	}

	result := make([][]byte, 0, len(values))
	for c := range p.ChildrenRange(1, 4) {
		v := make([]byte, 2)
		if _, err := raw.NewBufferFromSlice(c).Read(v); err != nil {
			t.Fatal(err)
		}

		result = append(result, v)
	}

	want := [][]byte{
		{3, 4},
		{5, 6},
		{7, 8}}

	if !slices.EqualFunc(result, want, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func createBuffer(t *testing.T, s *storage.Storage) *storage.Buffer {
	b, err := s.Request(t.Context(), link.BlockID{})
	if err != nil {
		t.Fatal(err)
	}

	return b
}
