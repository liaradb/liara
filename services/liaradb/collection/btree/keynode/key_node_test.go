package keynode

import (
	"slices"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
	"github.com/liaradb/liaradb/util/testing/testutil"
)

func TestKeyNode(t *testing.T) {
	t.Parallel()

	testutil.Run(t, "should insert", func(t *testing.T) {
		s := storagetesting.CreateStorage(t, 2, 256)
		b := createBuffer(t, s)
		defer b.Release()

		bp := node.New(b)
		kn := New(bp)

		data := testKeyNodeInsertData(t, kn)
		testKeyNodeChildren(t, kn, data)
	})

	t.Run("should replace root", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			s := storagetesting.CreateStorage(t, 2, 256)
			b := createBuffer(t, s)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := []keyEntry{
				{key.NewKey([]byte("key1")), link.FilePosition(1)},
				{key.NewKey([]byte("key2")), link.FilePosition(2)},
				{key.NewKey([]byte("key3")), link.FilePosition(3)},
			}

			for _, item := range data {
				if _, _, ok := kn.Insert(item.key, item.block); !ok {
					t.Error("should insert")
				}
			}

			testKeyNodeChildren(t, kn, data)

			want := []keyEntry{
				{key.NewKey([]byte("key1")), link.FilePosition(10)},
				{key.NewKey([]byte("key11")), link.FilePosition(11)},
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

	testutil.Run(t, "should iterate in order", func(t *testing.T) {
		s := storagetesting.CreateStorage(t, 2, 256)
		b := createBuffer(t, s)
		defer b.Release()

		bp := node.New(b)
		kn := New(bp)

		data := testKeyNodeInsertData(t, kn)
		testKeyNodeChildren(t, kn, data)
	})

	t.Run("should get child", func(t *testing.T) {
		t.Parallel()
		synctest.Test(t, func(t *testing.T) {
			s := storagetesting.CreateStorage(t, 2, 256)
			b := createBuffer(t, s)
			defer b.Release()

			bp := node.New(b)
			kn := New(bp)

			data := testKeyNodeInsertData(t, kn)
			testKeyNodeChild(t, kn, data)
		})
	})

	testutil.Run(t, "should search items", func(t *testing.T) {
		s := storagetesting.CreateStorage(t, 2, 256)
		b := createBuffer(t, s)
		defer b.Release()

		bp := node.New(b)
		kn := New(bp)

		data := testKeyNodeInsertData(t, kn)
		testKeyNodeSearch(t, kn, data)
	})

	// TODO: This method is private
	testutil.Run(t, "should search indexes", func(t *testing.T) {
		s := storagetesting.CreateStorage(t, 2, 256)
		b := createBuffer(t, s)
		defer b.Release()

		bp := node.New(b)
		kn := New(bp)

		data := testKeyNodeInsertData(t, kn)
		for i, e := range data {
			index := kn.searchIndex(e.key)
			if index != int16(i) {
				t.Errorf("incorrect index: %v, expected: %v", index, int16(i))
			}
		}
		{
			before := int16(0)
			result := kn.searchIndex(key.NewKey([]byte("a")))
			if result != before {
				t.Errorf("incorrect before: %v, expected: %v", result, before)
			}
		}
		{
			after := int16(len(data))
			result := kn.searchIndex(key.NewKey([]byte("e")))
			if result != after {
				t.Errorf("incorrect after: %v, expected: %v", result, after)
			}
		}
	})

	testutil.Run(t, "should search before items", func(t *testing.T) {
		s := storagetesting.CreateStorage(t, 2, 256)
		b := createBuffer(t, s)
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

	testutil.Run(t, "should search after items", func(t *testing.T) {
		s := storagetesting.CreateStorage(t, 2, 256)
		b := createBuffer(t, s)
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

	testutil.Run(t, "should fill", func(t *testing.T) {
		s := storagetesting.CreateStorage(t, 2, 256)
		b := createBuffer(t, s)
		defer b.Release()

		bp := node.New(b)
		kn := New(bp)

		data := []keyEntry{
			{key.NewKey([]byte("b")), 1},
			{key.NewKey([]byte("c")), 2},
			{key.NewKey([]byte("d")), 3},
		}

		kn.Fill(1, func(yield func(k key.Key, fp link.FilePosition) bool) {
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

	testutil.Run(t, "should replace", func(t *testing.T) {
		s := storagetesting.CreateStorage(t, 2, 256)
		b := createBuffer(t, s)
		defer b.Release()

		bp := node.New(b)
		kn := New(bp)

		data := []keyEntry{
			{key.NewKey([]byte("b")), 1},
			{key.NewKey([]byte("c")), 2},
			{key.NewKey([]byte("d")), 3},
		}

		kn.Replace(1, func(yield func(k key.Key, fp link.FilePosition) bool) {
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
		key, block, ok := kn.Child(i)
		if ok {
			result = append(result, newKeyEntry(key, block))
		}
	}

	if _, _, ok := kn.Child(kn.Count()); ok {
		t.Error("should not get beyond count")
	}

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
	b, err := s.Request(t.Context(), link.BlockID{})
	if err != nil {
		t.Fatal(err)
	}

	return b
}
