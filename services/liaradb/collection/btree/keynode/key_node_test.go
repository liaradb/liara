package keynode

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/transaction/record"
	"github.com/liaradb/liaradb/util/testing/logtesting"
	"github.com/liaradb/liaradb/util/testing/should"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestKeyNode(t *testing.T) {
	t.Parallel()

	t.Run("should insert", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := testKeyNodeInsertData(t, kn)
			testKeyNodeChildren(t, kn, data)
		})
	})

	t.Run("should insert and split", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 80, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			n := node.New(b)
			kn := New(n)

			data := []keyEntry{
				{key.NewKey([]byte("abcd")), 1},
				{key.NewKey([]byte("efgh")), 2},
				{key.NewKey([]byte("ijkl")), 3},
			}

			if a, b, ok := kn.Insert(data[0].key, data[0].block); !ok {
				t.Error("should insert")
			} else if a != nil {
				t.Errorf("incorrect iterator: %v, expected: %v", a, nil)
			} else if b != nil {
				t.Errorf("incorrect iterator: %v, expected: %v", b, nil)
			}

			if a, b, ok := kn.Insert(data[1].key, data[1].block); !ok {
				t.Error("should insert")
			} else if a != nil {
				t.Errorf("incorrect iterator: %v, expected: %v", a, nil)
			} else if b != nil {
				t.Errorf("incorrect iterator: %v, expected: %v", b, nil)
			}

			if a, b, ok := kn.Insert(data[2].key, data[2].block); ok {
				t.Error("should not insert")
			} else {
				resultA := make([]keyEntry, 0, 1)
				for k, v := range a {
					resultA = append(resultA, newKeyEntry(k, v))
				}
				if want := data[:1]; !slices.Equal(resultA, want) {
					t.Errorf("incorrect result: %v, expected: %v", resultA, want)
				}

				resultB := make([]keyEntry, 0, 2)
				for k, v := range b {
					resultB = append(resultB, newKeyEntry(k, v))
				}
				if want := data[1:]; !slices.Equal(resultB, want) {
					t.Errorf("incorrect result: %v, expected: %v", resultB, want)
				}
			}

			l := logtesting.MockLog{}
			if err := kn.LogActions(t.Context(), &l); err != nil {
				t.Error(err)
			}

			for i, r := range l.Records() {
				if a := r.Action(); a != record.ActionInsert {
					t.Errorf("incorrect action: %v, expected: %v", a, record.ActionInsert)
				}
				var ke keyEntry
				ke.Read(r.Data())
				switch i {
				case 0:
					if k, want := ke.Key(), data[0].Key(); k != want {
						t.Errorf("incorrect key: %v, expected: %v", k, want)
					}
					if b, want := ke.Block(), data[0].Block(); b != want {
						t.Errorf("incorrect key: %v, expected: %v", b, want)
					}
				case 1:
					if k, want := ke.Key(), data[1].Key(); k != want {
						t.Errorf("incorrect key: %v, expected: %v", k, want)
					}
					if b, want := ke.Block(), data[1].Block(); b != want {
						t.Errorf("incorrect key: %v, expected: %v", b, want)
					}
				default:
					t.Error("too many actions")
				}
			}
		})
	})

	t.Run("should replace root", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := []keyEntry{
				{key.NewKey([]byte("key1")), link.BlockID(1)},
				{key.NewKey([]byte("key2")), link.BlockID(2)},
				{key.NewKey([]byte("key3")), link.BlockID(3)},
			}

			for _, item := range data {
				if _, _, ok := kn.Insert(item.key, item.block); !ok {
					t.Error("should insert")
				}
			}

			testKeyNodeChildren(t, kn, data)

			want := []keyEntry{
				{key.NewKey([]byte("key1")), link.BlockID(10)},
				{key.NewKey([]byte("key11")), link.BlockID(11)},
			}

			if ok := kn.ReplaceRoot(1, want[0].block, want[1].key, want[1].block); !ok {
				t.Error("should replace")
			}

			testKeyNodeChildren(t, kn, want)

			if l := kn.Level(); l != 1 {
				t.Errorf("incorrect level: %v, expected: %v", l, 1)
			}
		})
	})

	t.Run("should iterate in order", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := testKeyNodeInsertData(t, kn)
			testKeyNodeChildren(t, kn, data)
		})
	})

	t.Run("should get child", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := testKeyNodeInsertData(t, kn)
			testKeyNodeChild(t, kn, data)
		})
	})

	t.Run("should search items", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := testKeyNodeInsertData(t, kn)
			testKeyNodeSearch(t, kn, data)
		})
	})

	// TODO: This method is private
	t.Run("should search indexes", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := testKeyNodeInsertData(t, kn)
			for i, e := range data {
				index := kn.searchIndex(e.key)
				if index != link.SlotID(i) {
					t.Errorf("incorrect index: %v, expected: %v", index, int16(i))
				}
			}
			{
				before := link.SlotID(0)
				result := kn.searchIndex(key.NewKey([]byte("a")))
				if result != before {
					t.Errorf("incorrect before: %v, expected: %v", result, before)
				}
			}
			{
				after := link.SlotID(len(data))
				result := kn.searchIndex(key.NewKey([]byte("e")))
				if result != after {
					t.Errorf("incorrect after: %v, expected: %v", result, after)
				}
			}
		})
	})

	t.Run("should search before items", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := testKeyNodeInsertData(t, kn)
			want := data[0].block
			result := kn.Search(key.NewKey([]byte("a")))
			if result != want {
				t.Errorf("incorrect result: %v, expected: %v", result, want)
			}
		})
	})

	t.Run("should search after items", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := testKeyNodeInsertData(t, kn)
			want := data[2].block
			result := kn.Search(key.NewKey([]byte("e")))
			if result != want {
				t.Errorf("incorrect result: %v, expected: %v", result, want)
			}
		})
	})

	t.Run("should fill", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := []keyEntry{
				{key.NewKey([]byte("b")), 1},
				{key.NewKey([]byte("c")), 2},
				{key.NewKey([]byte("d")), 3},
			}

			kn.Fill(1, func(yield func(k key.Key, fp link.BlockID) bool) {
				for _, ke := range data {
					if !yield(ke.Key(), ke.Block()) {
						return
					}
				}
			})

			testKeyNodeChildren(t, kn, data)

			if l := kn.Level(); l != 1 {
				t.Errorf("incorrect level: %v, expected: %v", l, 1)
			}
		})
	})

	t.Run("should replace", func(t *testing.T) {
		storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
			b := createBuffer(t, s.Storage)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := []keyEntry{
				{key.NewKey([]byte("b")), 1},
				{key.NewKey([]byte("c")), 2},
				{key.NewKey([]byte("d")), 3},
			}

			kn.Replace(1, func(yield func(k key.Key, fp link.BlockID) bool) {
				for _, ke := range data {
					if !yield(ke.Key(), ke.Block()) {
						return
					}
				}
			})

			testKeyNodeChildren(t, kn, data)

			if l := kn.Level(); l != 1 {
				t.Errorf("incorrect level: %v, expected: %v", l, 1)
			}
		})
	})
}

// Insert in mixed order
func testKeyNodeInsertData(t *testing.T, kn *KeyNode) []keyEntry {
	data := []keyEntry{
		{key.NewKey([]byte("b")), 1},
		{key.NewKey([]byte("c")), 2},
		{key.NewKey([]byte("d")), 3},
	}

	order := []int{0, 2, 1}
	for _, i := range order {
		e := data[i]
		if _, _, ok := kn.Insert(e.key, e.block); !ok {
			t.Error("should insert")
		}
	}

	return data
}

// Verify items are in order
func testKeyNodeChild(t *testing.T, kn *KeyNode, data []keyEntry) {
	result := make([]keyEntry, 0, len(data))
	for i := range kn.Count() {
		key, block := kn.Child(i)
		result = append(result, newKeyEntry(key, block))
	}

	should.Panic(t, "should not get beyond count", func() {
		_, _ = kn.Child(kn.Count())
	})

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

// Verify items are in order
func testKeyNodeChildren(t *testing.T, kn *KeyNode, data []keyEntry) {
	result := make([]keyEntry, 0, len(data))
	for key, block := range kn.Children() {
		result = append(result, newKeyEntry(key, block))
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

// Verify items are all searchable
func testKeyNodeSearch(t *testing.T, kn *KeyNode, data []keyEntry) {
	for _, e := range data {
		if block := kn.Search(e.key); block != e.block {
			t.Errorf("incorrect record id: %v, expected: %v", block, e.block)
		}
	}
}

func createBuffer(t *testing.T, s *storage.Storage) *storage.Buffer {
	b, err := s.Request(t.Context(), link.FileBlockID{})
	if err != nil {
		t.Fatal(err)
	}

	return b
}
