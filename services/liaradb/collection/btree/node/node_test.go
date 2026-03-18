package node

import (
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

const (
	testHeaderSize = 2 + headerSize
)

func TestNode_Dirty(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Dirty)
}

func testNode_Dirty(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	if b.Dirty() {
		t.Error("should not be dirty")
	}

	n := New(b)
	n.SetDirty()

	if !b.Dirty() {
		t.Error("should be dirty")
	}
}

func TestNode_Level(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Level)
}

func testNode_Level(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	n := New(b)

	if l := n.Level(); l != 0 {
		t.Errorf("incorrect level: %v, expected: %v", l, 0)
	}

	n.SetLevel(2)

	if l := n.Level(); l != 2 {
		t.Errorf("incorrect level: %v, expected: %v", l, 2)
	}
}

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

	n := New(b)
	v0 := []byte{1, 2, 3, 4, 5}
	v1 := []byte{6, 7, 8, 9, 10}

	if c := n.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	if s := n.Space(); s != s0 {
		t.Errorf("incorrect space: %v, expected: %v", s, s0)
	}

	i, b0, ok := n.Append(16)
	if !ok {
		t.Error("should get a buffer")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if c := n.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	if s := n.Space(); s != s1 {
		t.Errorf("incorrect space: %v, expected: %v", s, s1)
	}

	if _, err := buffer.NewFromSlice(b0).Write(v0); err != nil {
		t.Error(err)
	}

	i, b1, ok := n.Append(16)
	if !ok {
		t.Error("should get a buffer")
	} else if i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 1)
	}

	if c := n.Count(); c != 2 {
		t.Errorf("incorrect count: %v, expected: %v", c, 2)
	}

	if s := n.Space(); s != s2 {
		t.Errorf("incorrect space: %v, expected: %v", s, s2)
	}

	if _, err := buffer.NewFromSlice(b1).Write(v1); err != nil {
		t.Error(err)
	}

	// if _, err := raw.NewBufferFromSlice(b0).Seek(0, io.SeekStart); err != nil {
	// 	t.Error(err)
	// }

	r0 := make([]byte, 5)
	if _, err := buffer.NewFromSlice(b0).Read(r0); err != nil {
		t.Error(err)
	}

	if !slices.Equal(r0, v0) {
		t.Errorf("incorrect result: %v, expected: %v", r0, v0)
	}

	// if _, err := raw.NewBufferFromSlice(b1).Seek(0, io.SeekStart); err != nil {
	// 	t.Error(err)
	// }

	r1 := make([]byte, 5)
	if _, err := buffer.NewFromSlice(b1).Read(r1); err != nil {
		t.Error(err)
	}

	if !slices.Equal(r1, v1) {
		t.Errorf("incorrect result: %v, expected: %v", r1, v1)
	}

	synctest.Wait()
}

func TestNode_Insert(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Insert)
}

func testNode_Insert(t *testing.T) {
	const (
		size int16 = 256
		s0         = size - itemSize - testHeaderSize
		s1         = s0 - itemSize - 16
		s2         = s1 - itemSize - 16
	)

	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	n := New(b)
	v0 := []byte{1, 2, 3, 4, 5}
	v1 := []byte{6, 7, 8, 9, 10}

	if s := n.Space(); s != s0 {
		t.Fatalf("incorrect space: %v, expected: %v", s, s0)
	}

	i, b0, ok := n.Insert(16, 0)
	if !ok {
		t.Fatal("should get a buffer")
	} else if i != 0 {
		t.Fatalf("incorrect index: %v, expected: %v", i, 0)
	}

	if s := n.Space(); s != s1 {
		t.Fatalf("incorrect space: %v, expected: %v", s, s1)
	}

	if _, err := buffer.NewFromSlice(b0).Write(v0); err != nil {
		t.Fatal(err)
	}

	i, b1, ok := n.Insert(16, 1)
	if !ok {
		t.Fatal("should get a buffer")
	} else if i != 1 {
		t.Fatalf("incorrect index: %v, expected: %v", i, 1)
	}

	if s := n.Space(); s != s2 {
		t.Fatalf("incorrect space: %v, expected: %v", s, s2)
	}

	if _, err := buffer.NewFromSlice(b1).Write(v1); err != nil {
		t.Fatal(err)
	}

	// if _, err := raw.NewBufferFromSlice(b0).Seek(0, io.SeekStart); err != nil {
	// 	t.Fatal(err)
	// }

	r0 := make([]byte, 5)
	if _, err := buffer.NewFromSlice(b0).Read(r0); err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(r0, v0) {
		t.Fatalf("incorrect result: %v, expected: %v", r0, v0)
	}

	// if _, err := raw.NewBufferFromSlice(b1).Seek(0, io.SeekStart); err != nil {
	// 	t.Fatal(err)
	// }

	r1 := make([]byte, 5)
	if _, err := buffer.NewFromSlice(b1).Read(r1); err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(r1, v1) {
		t.Errorf("incorrect result: %v, expected: %v", r1, v1)
	}

	synctest.Wait()
}

func TestNode_Space(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Space)
}

func testNode_Space(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 16+itemSize+testHeaderSize)
	b := createBuffer(t, s)
	defer b.Release()

	n := New(b)

	if s := n.Space(); s != 16 {
		t.Errorf("incorrect space: %v, expected: %v", s, 16)
	}

	if _, _, ok := n.Append(16); !ok {
		t.Error("should get a buffer")
	}

	if s := n.Space(); s != 0 {
		t.Errorf("incorrect space: %v, expected: %v", s, 0)
	}

	if _, _, ok := n.Append(16); ok {
		t.Error("should not get a buffer")
	}

	synctest.Wait()
}

func TestNode_Child(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Child)
}

func testNode_Child(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	n := New(b)
	values := [][]byte{
		{1, 2, 3, 4, 5},
		{6, 7, 8, 9, 10}}
	_, b0, ok := n.Append(16)
	if !ok {
		t.Error("should get a buffer")
	}

	if _, err := buffer.NewFromSlice(b0).Write(values[0]); err != nil {
		t.Error(err)
	}

	_, b1, ok := n.Append(16)
	if !ok {
		t.Error("should get a buffer")
	}

	if _, err := buffer.NewFromSlice(b1).Write(values[1]); err != nil {
		t.Error(err)
	}

	result := make([][]byte, 0, 2)
	for i := range 2 {
		c, ok := n.Child(int16(i))
		if !ok {
			t.Fatal("should get a buffer")
		}

		v := make([]byte, 5)
		if _, err := buffer.NewFromSlice(c).Read(v); err != nil {
			t.Fatal(err)
		}

		result = append(result, v)
	}

	if !slices.EqualFunc(result, values, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, values)
	}

	synctest.Wait()
}

func TestNode_Children(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Children)
}

func testNode_Children(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	n := New(b)
	values := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10}}

	for _, v := range values {
		_, b, ok := n.Append(itemSize * int16(len(v)))
		if !ok {
			t.Error("should get a buffer")
		}

		if _, err := buffer.NewFromSlice(b).Write(v); err != nil {
			t.Error(err)
		}
	}

	result := make([][]byte, 0, len(values))
	for c := range n.Children() {
		v := make([]byte, 2)
		if _, err := buffer.NewFromSlice(c).Read(v); err != nil {
			t.Fatal(err)
		}

		result = append(result, v)
	}

	if !slices.EqualFunc(result, values, slices.Equal) {
		t.Errorf("incorrect result: %v, expected: %v", result, values)
	}

	synctest.Wait()
}

func TestNode_ChildrenRange(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_ChildrenRange)
}

func testNode_ChildrenRange(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)
	defer b.Release()

	n := New(b)
	values := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10}}

	for _, v := range values {
		_, b0, ok := n.Append(itemSize * int16(len(v)))
		if !ok {
			t.Error("should get a buffer")
		}

		if _, err := buffer.NewFromSlice(b0).Write(v); err != nil {
			t.Error(err)
		}
	}

	result := make([][]byte, 0, len(values))
	for c := range n.ChildrenRange(1, 4) {
		v := make([]byte, 2)
		if _, err := buffer.NewFromSlice(c).Read(v); err != nil {
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

	synctest.Wait()
}

func createBuffer(t *testing.T, s *storage.Storage) *storage.Buffer {
	b, err := s.Request(t.Context(), link.BlockID{})
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func TestNode_Clear(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Clear)
}

func testNode_Clear(t *testing.T) {
	s := storagetesting.CreateStorage(t, 2, 256)
	b := createBuffer(t, s)

	n := New(b)
	values := [][]byte{
		{1, 2},
		{3, 4},
		{5, 6},
		{7, 8},
		{9, 10}}

	for _, v := range values {
		_, b, ok := n.Append(itemSize * int16(len(v)))
		if !ok {
			t.Error("should get a buffer")
		}

		if _, err := buffer.NewFromSlice(b).Write(v); err != nil {
			t.Error(err)
		}
	}

	n.Clear()

	result := make([][]byte, 0, len(values))
	for c := range n.Children() {
		v := make([]byte, 2)
		if _, err := buffer.NewFromSlice(c).Read(v); err != nil {
			t.Fatal(err)
		}

		result = append(result, v)
	}

	if l := len(result); l != 0 {
		t.Errorf("incorrect length: %v, expected: %v", l, 0)
	}

	n.Release()

	synctest.Wait()
}
