package btree

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
	"github.com/liaradb/liaradb/util/testing/txtesting"
)

// TODO: Test latching

func TestCursor_GetRoot_Default(t *testing.T) {
	storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		fn := link.NewFileName("testfile")
		if l, err := NewCursor(tx).Level(ctx, fn); err != nil {
			t.Error(err)
		} else if l != 0 {
			t.Errorf("incorrect level: %v, expected: %v", l, 0)
		}
	})
}

func TestCursor_Insert__Root(t *testing.T) {
	storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		fn := link.NewFileName("testfile")

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		data := []leafEntry{
			newLeafEntry(
				key.NewKey([]byte("a")),
				link.NewRecordID(1, 2)),
			newLeafEntry(
				key.NewKey([]byte("b")),
				link.NewRecordID(3, 4)),
			newLeafEntry(
				key.NewKey([]byte("c")),
				link.NewRecordID(5, 6)),
		}

		for _, e := range data {
			if err := NewCursor(tx).Insert(ctx, tx, fn, e.key, e.recordID); err != nil {
				t.Error(err)
			}
		}

		for _, e := range data {
			if rid, err := NewCursor(tx).Search(ctx, fn, e.key); err != nil {
				t.Fatal(err)
			} else if rid != e.recordID {
				t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
			}
		}

		if l, err := NewCursor(tx).Level(ctx, fn); err != nil {
			t.Error(err)
		} else if l != 0 {
			t.Errorf("incorrect level: %v, expected: %v", l, 0)
		}
	})
}

func TestCursor_Insert__RootSplit(t *testing.T) {
	storagetesting.SyncTest(t, 40, 80, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		fn := link.NewFileName("testfile")

		data := createData()
		for _, e := range data {
			if err := NewCursor(tx).Insert(ctx, tx, fn, e.key, e.recordID); err != nil {
				t.Fatal(e.key, err)
			}
		}

		for _, e := range data {
			if rid, err := NewCursor(tx).Search(ctx, fn, e.key); err != nil {
				t.Error(err, e.key)
			} else if rid != e.recordID {
				t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
			}
		}

		if l, err := NewCursor(tx).Level(ctx, fn); err != nil {
			t.Error(err)
		} else if l <= 2 {
			t.Errorf("incorrect level: %v, expected: %v", "> 2", l)
		}
	})
}

func TestCursor_Insert__Reverse(t *testing.T) {
	storagetesting.SyncTest(t, 40, 80, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		fn := link.NewFileName("testfile")

		data := createData()
		for _, e := range reverseData(data) {
			if err := NewCursor(tx).Insert(ctx, tx, fn, e.key, e.recordID); err != nil {
				t.Fatal(e.key, err)
			}
		}

		for _, e := range data {
			if rid, err := NewCursor(tx).Search(ctx, fn, e.key); err != nil {
				t.Error(err, e.key)
			} else if rid != e.recordID {
				t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
			}
		}
	})
}

func TestCursor_Insert__Random(t *testing.T) {
	storagetesting.SyncTest(t, 40, 80, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		fn := link.NewFileName("testfile")

		data := createData()

		// Insert in mixed order
		order := []int{
			5,
			8,
			2,
			4,
			3,
			6,
			0,
			7,
			1,
		}

		for i, e := range reorderData(order, data) {
			if err := NewCursor(tx).Insert(ctx, tx, fn, e.key, e.recordID); err != nil {
				t.Fatal(i, e.key, err)
			}
		}

		for _, i := range order {
			e := data[i]
			if rid, err := NewCursor(tx).Search(ctx, fn, e.key); err != nil {
				t.Error(err, e.key)
			} else if rid != e.recordID {
				t.Errorf("incorrect record id: %v, expected: %v", rid, e.recordID)
			}
		}
	})
}

func TestCursor_Insert__Existing(t *testing.T) {
	storagetesting.SyncTest(t, 8, 62, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		fn := link.NewFileName("testfile")

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		le := newLeafEntry(
			key.NewKey([]byte("0")),
			link.NewRecordID(1, 2))

		if err := NewCursor(tx).Insert(ctx, tx, fn, le.key, le.recordID); err != nil {
			t.Fatal(le.key, err)
		}

		if err := NewCursor(tx).Insert(ctx, tx, fn, le.key, le.recordID); err == nil {
			t.Error("should not insert the same key")
		}

		c := 0
		for _, err := range NewCursor(tx).All(ctx, fn, 0, 0) {
			if err != nil {
				t.Error(err)
			}

			c++
		}

		if c != 1 {
			t.Errorf("incorrect count: %v, expected: %v", c, 1)
		}
	})
}

func TestCursor_SearchRange(t *testing.T) {
	storagetesting.SyncTest(t, 40, 80, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		fn := link.NewFileName("testfile")

		data := createData()

		for _, e := range data {
			if err := NewCursor(tx).Insert(ctx, tx, fn, e.key, e.recordID); err != nil {
				t.Fatal(e.key, err)
			}
		}

		wantAll := make([]SearchResult, 0, len(data))
		for _, e := range data {
			wantAll = append(wantAll, SearchResult{e.key, e.recordID})
		}

		for i, e := range data {
			c := NewCursor(tx)
			result := make([]SearchResult, 0, len(data))
			for rid, err := range c.SearchRange(ctx, fn, e.key, 0, 0) {
				if err != nil {
					t.Fatal(err)
				}

				result = append(result, rid)
			}

			want := wantAll[i:]
			if !slices.Equal(result, want) {
				t.Errorf("incorrect result: %v, expected: %v", result, want)
			}
		}

		// Skip and Limit
		{
			c := NewCursor(tx)
			result := make([]SearchResult, 0, len(data))
			for r, err := range c.SearchRange(ctx, fn, key.NewKey([]byte("1")), 1, 3) {
				if err != nil {
					t.Fatal(err)
				}

				result = append(result, r)
			}

			want := wantAll[2:5]
			if !slices.Equal(result, want) {
				t.Errorf("incorrect result: %v, expected: %v", result, want)
			}
		}
	})
}

func TestCursor_All(t *testing.T) {
	storagetesting.SyncTest(t, 40, 80, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		fn := link.NewFileName("testfile")

		data := createData()

		for _, e := range data {
			if err := NewCursor(tx).Insert(ctx, tx, fn, e.key, e.recordID); err != nil {
				t.Fatal(e.key, err)
			}
		}

		wantAll := make([]link.RecordID, 0, len(data))
		for _, e := range data {
			wantAll = append(wantAll, e.recordID)
		}

		for i := range data {
			c := NewCursor(tx)
			result := make([]link.RecordID, 0, len(data))
			for r, err := range c.All(ctx, fn, i, 0) {
				if err != nil {
					t.Fatal(err)
				}

				result = append(result, r.RecordID)
			}

			want := wantAll[i:]
			if !slices.Equal(result, want) {
				t.Errorf("incorrect result: %v, expected: %v", result, want)
			}
		}

		// Skip and Limit
		{
			c := NewCursor(tx)
			result := make([]link.RecordID, 0, len(data))
			for r, err := range c.All(ctx, fn, 1, 3) {
				if err != nil {
					t.Fatal(err)
				}

				result = append(result, r.RecordID)
			}

			want := wantAll[1:4]
			if !slices.Equal(result, want) {
				t.Errorf("incorrect result: %v, expected: %v", result, want)
			}
		}
	})
}

func reverseData(data []leafEntry) []leafEntry {
	data = slices.Clone(data)
	slices.Reverse(data)
	return data
}

func reorderData(order []int, data []leafEntry) []leafEntry {
	data = slices.Clone(data)

	hash := make(map[key.Key]int)
	for i, le := range data {
		if i >= len(order) {
			break
		}

		hash[le.key] = order[i]
	}
	l := len(data)

	slices.SortFunc(data, func(a, b leafEntry) int {
		c, ok := hash[a.key]
		if !ok {
			c = l
		}

		d, ok := hash[b.key]
		if !ok {
			d = l
		}

		return c - d
	})

	return data
}

type leafEntry struct {
	key      key.Key
	recordID link.RecordID
}

func newLeafEntry(key key.Key, recordID link.RecordID) leafEntry {
	return leafEntry{
		key:      key,
		recordID: recordID,
	}
}

func createData() []leafEntry {
	return []leafEntry{
		newLeafEntry(
			key.NewKey([]byte("0")),
			link.NewRecordID(1, 2)),
		newLeafEntry(
			key.NewKey([]byte("1")),
			link.NewRecordID(3, 4)),
		newLeafEntry(
			key.NewKey([]byte("2")),
			link.NewRecordID(5, 6)),
		newLeafEntry(
			key.NewKey([]byte("3")),
			link.NewRecordID(7, 8)),
		newLeafEntry(
			key.NewKey([]byte("4")),
			link.NewRecordID(9, 10)),
		newLeafEntry(
			key.NewKey([]byte("5")),
			link.NewRecordID(11, 12)),
		newLeafEntry(
			key.NewKey([]byte("6")),
			link.NewRecordID(13, 14)),
		newLeafEntry(
			key.NewKey([]byte("7")),
			link.NewRecordID(15, 16)),
		newLeafEntry(
			key.NewKey([]byte("8")),
			link.NewRecordID(17, 18)),
	}
}
