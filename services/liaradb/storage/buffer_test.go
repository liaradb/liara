package storage

import (
	"io"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/encoder/buffer"
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

func TestBuffer_ReadWrite(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testBuffer_ReadWrite)
}

func testBuffer_ReadWrite(t *testing.T) {
	s := createStorage(t, 2, 16)
	ctx := t.Context()

	b, err := s.Request(ctx, link.NewBlockID(link.NewFileName(""), 0))
	if err != nil {
		t.Fatal(err)
	}

	if c := b.Cursor(); c != 0 {
		t.Errorf("incorrect cursor: %v, expected: %v", c, 0)
	}

	data := []byte{1, 2, 3}

	if l, err := b.Write(data); err != nil {
		t.Fatal(err)
	} else if l != len(data) {
		t.Fatalf("incorrect length: %v, expected: %v", l, 3)
	}

	if l, err := b.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	} else if l != 0 {
		t.Fatalf("incorrect length: %v, expected: %v", l, 0)
	}

	result := make([]byte, 3)
	if l, err := b.Read(result); err != nil {
		t.Fatal(err)
	} else if l != len(data) {
		t.Fatalf("incorrect length: %v, expected: %v", l, 3)
	}

	if !slices.Equal(result, data) {
		t.Error("Should copy")
	}

	b.Release()

	synctest.Wait()
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

	if _, err := b0.Write([]byte{1, 2, 3}); err != nil {
		t.Fatal(err)
	}

	b1.Clone(b0)

	if !slices.Equal(b0.Raw(), b1.Raw()) {
		t.Error("Should copy")
	}

	b0.Release()
	b1.Release()

	synctest.Wait()
}

func TestBuffer_Seek(t *testing.T) {
	t.Parallel()

	const initialPosition = 10

	for message, c := range map[string]struct {
		skip     bool
		position int64
		whence   int
		err      error
		n        int
	}{
		"should handle defaults": {
			position: 0,
			whence:   0,
			err:      nil,
			n:        0},
		"should not seek to negative position from start": {
			position: -1,
			whence:   io.SeekStart,
			err:      buffer.ErrUnderflow,
			n:        initialPosition},
		"should not seek to negative position from default": {
			position: -1,
			whence:   10,
			err:      buffer.ErrUnderflow,
			n:        initialPosition},
		"should not seek to negative position from current": {
			position: -11,
			whence:   io.SeekCurrent,
			err:      buffer.ErrUnderflow,
			n:        initialPosition},
		"should not seek to negative position from end": {
			position: -21,
			whence:   io.SeekEnd,
			err:      buffer.ErrUnderflow,
			n:        initialPosition},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			synctest.Test(t, func(t *testing.T) {
				s := createStorage(t, 1, 20)
				ctx := t.Context()

				b, err := s.Request(ctx, link.NewBlockID(link.NewFileName(""), 0))
				if err != nil {
					t.Fatal(err)
				}

				if n, err := b.Seek(initialPosition, io.SeekStart); err != nil {
					t.Error(err)
				} else if n != initialPosition {
					t.Errorf("%v: incorrect count: %v, expected: %v", message, n, initialPosition)
				}

				n, err := b.Seek(c.position, c.whence)
				if err != c.err {
					if c.err == nil {
						t.Error(err)
					} else {
						t.Errorf("%v: incorrect error: %v, expected: %v", message, err, c.err)
					}
				}
				if n != int64(c.n) {
					t.Errorf("%v: incorrect n: %v, expected: %v", message, n, c.n)
				}

				if o := b.Cursor(); o != int64(c.n) {
					t.Errorf("%v: incorrect cursor: %v, expected: %v", message, o, c.n)
				}

				b.Release()

				synctest.Wait()
			})
		})
	}
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
	if _, err := b.Write(data); err != nil {
		t.Fatal(err)
	}

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
