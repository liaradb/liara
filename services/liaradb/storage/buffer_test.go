package storage

import (
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/storage/link"
)

func TestBuffer_Latch(t *testing.T) {
	t.Parallel()
	t.Skip()
	synctest.Test(t, testBuffer_Latch)
}

func testBuffer_Latch(t *testing.T) {
	// b := Buffer{}
	value := 0

	go func() {
		// b.Latch()
		// defer b.Unlatch()
		value0 := value
		time.Sleep(1 * time.Second)
		value = value0 + 1
	}()

	go func() {
		// time.Sleep(1 * time.Second)
		value1 := value
		time.Sleep(1 * time.Second)
		value = value1 + 1
		// b.Latch()
		// defer b.Unlatch()
	}()

	time.Sleep(10 * time.Second)
	if value != 2 {
		t.Errorf("incorrect value: %v, expected: %v", value, 2)
	}
}

func TestBuffer_Clone(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testBuffer_Clone)
}

func testBuffer_Clone(t *testing.T) {
	s := createStorage(t, 2, 16)
	ctx := t.Context()

	b0, err := s.Request(ctx, link.NewBlockID(link.NewFileName(""), 0))
	if err != nil {
		t.Fatal(err)
	}

	b1, err := s.Request(ctx, link.NewBlockID(link.NewFileName(""), 0))
	if err != nil {
		t.Fatal(err)
	}

	if r := b0.Reads(); r != 2 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 2)
	}

	if r := b1.Reads(); r != 2 {
		t.Errorf("incorrect reads: %v, expected: %v", r, 2)
	}

	b0.Fill([]byte{1, 2, 3})

	b1.Clone(b0)

	if !slices.Equal(b0.Raw(), b1.Raw()) {
		t.Error("Should copy")
	}

	b0.Release()
	b1.Release()

	synctest.Wait()
}

func TestNode_Clear(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testNode_Clear)
}

func testNode_Clear(t *testing.T) {
	s := createStorage(t, 2, 8)
	b := createBuffer(t, s)

	base := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	empty := make([]byte, 8)

	data := slices.Clone(base)
	b.Fill(base)

	if !b.Dirty() {
		t.Error("should be dirty")
	}

	if s := b.Status(); s != BufferStatusDirty {
		t.Errorf("incorrect status: %v, expected: %v", s, BufferStatusDirty)
	}

	if !slices.Equal(data, base) {
		t.Error("should not change data")
	}

	b.Clear()

	if s := b.Status(); s != BufferStatusUninitialized {
		t.Errorf("incorrect status: %v, expected: %v", s, BufferStatusUninitialized)
	}

	if !slices.Equal(data, base) {
		t.Error("should not change data")
	}

	if raw := b.Raw(); !slices.Equal(raw, empty) {
		t.Errorf("incorrect data: %v, expected: %v", raw, empty)
	}

	b.Release()

	synctest.Wait()
}

func createBuffer(t *testing.T, s *Storage) *Buffer {
	b, err := s.Request(t.Context(), link.BlockID{})
	if err != nil {
		t.Fatal(err)
	}

	return b
}
