package heap

import (
	"context"
	"testing"

	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/should"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
	"github.com/liaradb/liaradb/util/testing/txtesting"
)

func TestHeap_Get(t *testing.T) {
	storagetesting.SyncTest(t, 100, 128, func(t *testing.T, s storagetesting.Storage) {
		h := New()

		fn := link.NewFileName("test")
		bid := link.BlockID(0)
		rid := link.NewRecordID(bid, link.SlotID(0))

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		should.Panic(t, "should not get value", func() {
			_, _ = h.Get(t.Context(), tx, fn, rid)
		})

		// Long name
		d := make([]byte, 200)
		copy(d, "name")
		want := NewTestItem(NewTestName(string(d)))
		rid, err := h.Insert(t.Context(), tx, fn, want)
		if err != nil {
			t.Fatal(err)
		}

		r, err := h.Get(t.Context(), tx, fn, rid)
		if err != nil {
			t.Fatal(err)
		}

		defer r.Release()

		var item TestItem
		if err := item.Read(r); err != nil {
			t.Fatal(err)
		}

		if item != *want {
			t.Errorf("incorrect result: %v, expected: %v", item, *want)
		}
	})
}

func TestHeap_Get__Error(t *testing.T) {
	storagetesting.SyncTest(t, 100, 128, func(t *testing.T, s storagetesting.Storage) {
		h := New()

		fn := link.NewFileName("test")
		bid := link.BlockID(0)
		rid := link.NewRecordID(bid, link.SlotID(0))

		tx := txtesting.New(s.Storage)

		// Long name
		d := make([]byte, 200)
		copy(d, "name")
		want := NewTestItem(NewTestName(string(d)))
		rid, err := h.Insert(t.Context(), tx, fn, want)
		if err != nil {
			t.Fatal(err)
		}

		r, err := h.Get(t.Context(), tx, fn, rid)
		if err != nil {
			t.Fatal(err)
		}

		r.Release()
		tx.Release()

		tx = txtesting.New(s.Storage)
		defer tx.Release()

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		_, err = h.Get(ctx, tx, fn, rid)
		if err == nil {
			t.Error("should return error")
		}
	})
}

func TestHeap_Delete(t *testing.T) {
	storagetesting.SyncTest(t, 100, 1024, func(t *testing.T, s storagetesting.Storage) {
		h := New()
		fn := link.NewFileName("test")
		want := NewTestItem(NewTestName("name"))

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		rid, err := h.Insert(t.Context(), tx, fn, want)
		if err != nil {
			t.Fatal(err)
		}

		if err := h.Delete(t.Context(), tx, fn, rid); err != nil {
			t.Fatal(err)
		}

		should.Panic(t, "should not get value", func() {
			_, _ = h.Get(t.Context(), tx, fn, rid)
		})
	})
}

func TestHeap_Replace(t *testing.T) {
	fn := link.NewFileName("test")

	t.Run("should handle same size", func(t *testing.T) {
		// Very small to span fragments
		storagetesting.SyncTest(t, 100, 40, func(t *testing.T, s storagetesting.Storage) {
			h := New()

			tx := txtesting.New(s.Storage)
			defer tx.Release()

			first := NewTestItem(NewTestName("name0"))
			rid, err := h.Insert(t.Context(), tx, fn, first)
			if err != nil {
				t.Fatal(err)
			}

			second := NewTestItem(NewTestName("name1"))
			if err := h.Replace(t.Context(), tx, fn, rid, second); err != nil {
				t.Fatal(err)
			}

			r, err := h.Get(t.Context(), tx, fn, rid)
			if err != nil {
				t.Fatal(err)
			}

			defer r.Release()

			var result TestItem
			if err := result.Read(r); err != nil {
				t.Fatal(err)
			}

			if result != *second {
				t.Errorf("incorrect result: %v, expected: %v", result, *second)
			}
		})
	})

	t.Run("should handle smaller", func(t *testing.T) {
		t.Skip()
		storagetesting.SyncTest(t, 100, 1024, func(t *testing.T, s storagetesting.Storage) {
			h := New()

			tx := txtesting.New(s.Storage)
			defer tx.Release()

			first := NewTestItem(NewTestName("name0"))
			rid, err := h.Insert(t.Context(), tx, fn, first)
			if err != nil {
				t.Fatal(err)
			}

			second := NewTestItem(NewTestName("name"))
			if err := h.Replace(t.Context(), tx, fn, rid, second); err != nil {
				t.Fatal(err)
			}

			r, err := h.Get(t.Context(), tx, fn, rid)
			if err != nil {
				t.Fatal(err)
			}

			defer r.Release()

			var result TestItem
			if err := result.Read(r); err != nil {
				t.Fatal(err)
			}

			if result != *second {
				t.Errorf("incorrect result: %v, expected: %v", result, *second)
			}
		})
	})

	t.Run("should handle larger", func(t *testing.T) {
		t.Skip()
		storagetesting.SyncTest(t, 100, 1024, func(t *testing.T, s storagetesting.Storage) {
			h := New()

			tx := txtesting.New(s.Storage)
			defer tx.Release()

			first := NewTestItem(NewTestName("name0"))
			rid, err := h.Insert(t.Context(), tx, fn, first)
			if err != nil {
				t.Fatal(err)
			}

			second := NewTestItem(NewTestName("name12"))
			if err := h.Replace(t.Context(), tx, fn, rid, second); err != nil {
				t.Fatal(err)
			}

			r, err := h.Get(t.Context(), tx, fn, rid)
			if err != nil {
				t.Fatal(err)
			}

			defer r.Release()

			var result TestItem
			if err := result.Read(r); err != nil {
				t.Fatal(err)
			}

			if result != *second {
				t.Errorf("incorrect result: %v, expected: %v", result, *second)
			}
		})
	})
}
