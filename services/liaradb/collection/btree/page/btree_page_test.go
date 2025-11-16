package page

import (
	"io"
	"slices"
	"testing"
)

const (
	headerSize = 8 + BTreePageHeaderSize
	itemSize   = 4
)

func TestBTreePage(t *testing.T) {
	t.Parallel()

	const (
		size int32 = 256
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

	if _, err := b0.Write(v0); err != nil {
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

	if _, err := b1.Write(v1); err != nil {
		t.Error(err)
	}

	if _, err := b0.Seek(0, io.SeekStart); err != nil {
		t.Error(err)
	}

	r0 := make([]byte, 5)
	if _, err := b0.Read(r0); err != nil {
		t.Error(err)
	}

	if !slices.Equal(r0, v0) {
		t.Errorf("incorrect result: %v, expected: %v", r0, v0)
	}

	if _, err := b1.Seek(0, io.SeekStart); err != nil {
		t.Error(err)
	}

	r1 := make([]byte, 5)
	if _, err := b1.Read(r1); err != nil {
		t.Error(err)
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

func TestBTreePage_Level(t *testing.T) {
	t.Parallel()

	p := New(make([]byte, headerSize))
	if l := p.Level(); l != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.SetLevel(1)
	if l := p.Level(); l != 1 {
		t.Errorf("incorrect value: %v, expected: %v", p, 1)
	}
}

func TestBTreePage_ParentID(t *testing.T) {
	t.Parallel()

	p := New(make([]byte, headerSize))
	if id := p.ParentID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.SetParentID(1)
	if id := p.ParentID(); id != 1 {
		t.Errorf("incorrect value: %v, expected: %v", p, 1)
	}
}

func TestBTreePage_PrevID(t *testing.T) {
	t.Parallel()

	p := New(make([]byte, headerSize))
	if id := p.PrevID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.SetPrevID(1)
	if id := p.PrevID(); id != 1 {
		t.Errorf("incorrect value: %v, expected: %v", p, 1)
	}
}

func TestBTreePage_NextID(t *testing.T) {
	t.Parallel()

	p := New(make([]byte, headerSize))
	if id := p.NextID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.SetNextID(1)
	if id := p.NextID(); id != 1 {
		t.Errorf("incorrect value: %v, expected: %v", p, 1)
	}
}

func TestBTreePage_LowID(t *testing.T) {
	t.Parallel()

	p := New(make([]byte, headerSize))
	if id := p.LowID(); id != 0 {
		t.Errorf("incorrect value: %v, expected: %v", p, 0)
	}

	p.SetLowID(1)
	if id := p.LowID(); id != 1 {
		t.Errorf("incorrect value: %v, expected: %v", p, 1)
	}
}
