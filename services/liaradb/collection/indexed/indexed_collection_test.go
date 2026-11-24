package indexed

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"uuid"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/heap"
	"github.com/liaradb/liaradb/collection/heap/record"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
	"github.com/liaradb/liaradb/util/testing/txtesting"
)

func TestIndexedCollection_InsertAndGet(t *testing.T) {
	storagetesting.SyncTest(t, 16, 1024, func(t *testing.T, s storagetesting.Storage) {
		ic := New(heap.New())
		tn := &testTableName{}

		k := key.NewKey([]byte("abcde"))
		want := []byte{1, 2, 3, 4, 5}
		b := base.NewFixedBytes(want)

		tx := txtesting.New(s.Storage)

		if _, err := ic.Insert(t.Context(), tx, tn, k, b); err != nil {
			t.Fatal(err)
		}

		if !tx.IsXLocked(tn, k.String()) {
			t.Fatal("should be x locked")
		}

		tx.Release()

		tx = txtesting.New(s.Storage)

		r, err := ic.GetByKey(t.Context(), tx, tn, k)
		if err != nil {
			t.Fatal(err)
		}

		if !tx.IsSLocked(tn, k.String()) {
			t.Fatal("should be x locked")
		}

		tx.Release()

		b0 := base.NewFixedBytesFromSize(5)
		if err := b0.Read(r); err != nil {
			t.Fatal(err)
		}

		r.Release()
		result := b0.Value()

		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	})
}

func TestIndexedCollection_Insert(t *testing.T) {
	storagetesting.SyncTest(t, 100, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		ic := New(heap.New())
		tn := &testTableName{}

		data := createData()

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		if err := insertData(ctx, tx, ic, tn, data); err != nil {
			t.Fatal(err)
		}

		testGet(ctx, t, tx, ic, tn, data)
		testList(ctx, t, tx, data, ic, tn)
	})
}

func TestIndexedCollection_Insert__Reverse(t *testing.T) {
	storagetesting.SyncTest(t, 100, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		ic := New(heap.New())
		tn := &testTableName{}

		data := createData()
		slices.Reverse(data)

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		if err := insertData(ctx, tx, ic, tn, data); err != nil {
			t.Fatal(err)
		}

		testGet(ctx, t, tx, ic, tn, data)
		testList(ctx, t, tx, data, ic, tn)
	})
}

func TestIndexedCollection_Replace(t *testing.T) {
	storagetesting.SyncTest(t, 4, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		ic := New(heap.New())
		tn := &testTableName{}

		id := NewTestID()
		ti0 := NewTestItem(id, NewTestName("name0"))
		k := key.NewKey(id.Bytes())

		tx := txtesting.New(s.Storage)
		if _, err := ic.Insert(ctx, tx, tn, k, ti0); err != nil {
			t.Fatal(err)
		}

		if !tx.IsXLocked(tn, k.String()) {
			t.Fatal("should be x locked")
		}

		ti1 := NewTestItem(id, NewTestName("name1"))

		tx.Release()

		tx = txtesting.New(s.Storage)
		if err := ic.Replace(ctx, tx, tn, k, ti1); err != nil {
			t.Fatal(err)
		}

		if !tx.IsXLocked(tn, k.String()) {
			t.Fatal("should be x locked")
		}

		tx.Release()

		tx = txtesting.New(s.Storage)
		r, err := ic.GetByKey(ctx, tx, tn, k)
		if err != nil {
			t.Fatal(err)
		}

		if !tx.IsSLocked(tn, k.String()) {
			t.Fatal("should be x locked")
		}

		defer r.Release()

		tx.Release()

		var result TestItem
		if err := result.Read(r); err != nil {
			t.Fatal(err)
		}

		if result != *ti1 {
			t.Errorf("incorrect result: %v, expected: %v", result, *ti1)
		}
	})
}

type item struct {
	key   string
	value *TestItem
}

func createData() []item {
	count := 9
	data := uuid.UUID{}
	items := make([]item, 0, count)
	for i := range count {
		data[15] = byte(i) + 1
		rid, _ := NewTestIDFromString(data.String())
		items = append(items, item{fmt.Sprintf("%v", i+1), NewTestItem(rid, NewTestName(""))})
	}
	return items
}

func insertData(ctx context.Context, tx Transaction, ic *IndexedCollection, tn tablename.TableName, data []item) error {
	for _, i := range data {
		k := key.NewKey(i.value.ID().Bytes())
		if _, err := ic.Insert(ctx, tx, tn, k, i.value); err != nil {
			return err
		}
	}
	return nil
}

func testGet(
	ctx context.Context,
	t *testing.T,
	tx Transaction,
	ic *IndexedCollection,
	tn tablename.TableName,
	data []item,
) {
	t.Helper()

	for _, i := range data {
		k := key.NewKey(i.value.ID().Bytes())
		r, err := ic.GetByKey(ctx, tx, tn, k)
		if err != nil {
			t.Fatal(i.key, err)
		}

		rl, err := readTestItem(r)
		if err != nil {
			t.Fatal(i.key, err)
		}

		if *rl != *i.value {
			t.Errorf("incorrect result: %v, expected: %v", *rl, *i.value)
		}
	}
}

func testList(
	ctx context.Context,
	t *testing.T,
	tx Transaction,
	data []item,
	ic *IndexedCollection,
	tn tablename.TableName,
) {
	result, err := getListValues(ctx, tx, data, ic, tn)
	if err != nil {
		t.Fatal(err)
	}

	want := createSortedValues(data)
	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func getListValues(
	ctx context.Context,
	tx Transaction,
	data []item,
	ic *IndexedCollection,
	tn tablename.TableName,
) ([]TestItem, error) {
	result := make([]TestItem, 0, len(data))
	i := 0
	for r, err := range ic.List(ctx, tx, tn) {
		if err != nil {
			return nil, err
		}

		rl, err := readTestItem(r)
		if err != nil {
			return nil, err
		}

		result = append(result, *rl)
		i++
	}
	return result, nil
}

func readTestItem(r *record.Record) (*TestItem, error) {
	defer r.Release()

	rl := TestItem{}
	return &rl, rl.Read(r)
}

func createSortedValues(data []item) []TestItem {
	type tuple struct {
		key   key.Key
		value *TestItem
	}

	tuples := make([]tuple, 0, len(data))
	for _, i := range data {
		tuples = append(tuples, tuple{key.NewKey(i.value.ID().Bytes()), i.value})
	}
	slices.SortFunc(tuples, func(a, b tuple) int {
		return strings.Compare(a.key.String(), b.key.String())
	})
	want := make([]TestItem, 0, len(data))
	for _, t := range tuples {
		want = append(want, *t.value)
	}
	return want
}

type testTableName struct {
}

func (t *testTableName) Collection() link.FileName {
	return link.NewFileName("table.col")
}

func (t *testTableName) Index(i base.Uint32) link.FileName {
	return link.NewFileName(fmt.Sprintf("table--%v.bt", i))
}

func (t *testTableName) String() string {
	return "table"
}
