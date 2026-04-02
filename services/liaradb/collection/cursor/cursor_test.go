package cursor

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestCursor(t *testing.T) {
	storagetesting.SyncTest(t, 3, 256, func(t *testing.T, st storagetesting.Storage) {
		s := st.Storage
		ctx := t.Context()
		fn := link.NewFileName("testfile")

		b0, err := s.Request(ctx, fn.BlockID(0))
		if err != nil {
			t.Fatal(err)
		}

		b1, err := s.Request(ctx, fn.BlockID(1))
		if err != nil {
			t.Fatal(err)
		}

		b2, err := s.Request(ctx, fn.BlockID(2))
		if err != nil {
			t.Fatal(err)
		}

		c := New(b0, b1, b2)

		c.Release()
	})
}

func TestCursor_Writer(t *testing.T) {
	storagetesting.SyncTest(t, 3, 8, func(t *testing.T, st storagetesting.Storage) {
		s := st.Storage
		ctx := t.Context()
		fn := link.NewFileName("testfile")

		b0, err := s.Request(ctx, fn.BlockID(0))
		if err != nil {
			t.Fatal(err)
		}

		b1, err := s.Request(ctx, fn.BlockID(1))
		if err != nil {
			t.Fatal(err)
		}

		b2, err := s.Request(ctx, fn.BlockID(2))
		if err != nil {
			t.Fatal(err)
		}

		c := New(b0, b1, b2)
		defer c.Release()

		want := make([]byte, 0, 24)
		for i := range cap(want) {
			want = append(want, byte(i))
		}

		w := c.Writer()
		if _, err := w.Write(want); err != nil {
			t.Fatal(err)
		}

		result := make([]byte, 24)
		if _, err := c.Reader().Read(result); err != nil {
			t.Fatal(err)
		}

		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	})
}
