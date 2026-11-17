package indexed

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"testing/synctest"
	"uuid"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/heap"
	"github.com/liaradb/liaradb/collection/heap/record"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/logtesting"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestIndexedCollection_InsertAndGet(t *testing.T) {
	storagetesting.SyncTest(t, 16, 1024, func(t *testing.T, s storagetesting.Storage) {
		ic := New(heap.New(s.Storage), btree.NewCursor(s.Storage))
		tn := &testTableName{}

		k := key.NewKey([]byte("abcde"))
		want := []byte{1, 2, 3, 4, 5}
		b := base.NewFixedBytes(want)

		if _, err := ic.Insert(t.Context(), &logtesting.MockLog{}, tn, k, b); err != nil {
			t.Fatal(err)
		}

		r, err := ic.GetByKey(t.Context(), &logtesting.MockLog{}, tn, k)
		if err != nil {
			t.Fatal(err)
		}

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

func TestIndexedCollection(t *testing.T) {
	storagetesting.SyncTest(t, 20, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		ic := New(heap.New(s.Storage), btree.NewCursor(s.Storage))
		tn := &testTableName{}

		data := createData()
		slices.Reverse(data)

		if err := insertData(ctx, &logtesting.MockLog{}, ic, tn, data); err != nil {
			t.Fatal(err)
		}

		testGet(ctx, t, ic, tn, data)
		testList(ctx, t, data, ic, tn)

		synctest.Wait()
	})
}

func TestIndexedCollection__LargeBuffer(t *testing.T) {
	t.Skip()
	storagetesting.SyncTest(t, 2, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		ic := New(heap.New(s.Storage), btree.NewCursor(s.Storage))
		tn := &testTableName{}

		data := createData()

		if err := insertData(ctx, &logtesting.MockLog{}, ic, tn, data); err != nil {
			t.Fatal(err)
		}

		testGet(ctx, t, ic, tn, data)
		testList(ctx, t, data, ic, tn)

		synctest.Wait()
	})
}

func TestIndexedCollection_Replace(t *testing.T) {
	storagetesting.SyncTest(t, 4, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		ic := New(heap.New(s.Storage), btree.NewCursor(s.Storage))
		tn := &testTableName{}
		l := &logtesting.MockLog{}

		id := NewTestID()
		ti0 := NewTestItem(id, NewTestName("name0"))

		k := key.NewKey(id.Bytes())
		if _, err := ic.Insert(ctx, l, tn, k, ti0); err != nil {
			t.Fatal(err)
		}

		ti1 := NewTestItem(id, NewTestName("name1"))

		if err := ic.Replace(ctx, l, tn, k, ti1); err != nil {
			t.Fatal(err)
		}

		r, err := ic.GetByKey(ctx, l, tn, k)
		if err != nil {
			t.Fatal(err)
		}

		defer r.Release()

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
	ic *IndexedCollection,
	tn tablename.TableName,
	data []item,
) {
	t.Helper()

	for _, i := range data {
		k := key.NewKey(i.value.ID().Bytes())
		value, err := ic.GetByKey(ctx, &logtesting.MockLog{}, tn, k)
		if err != nil {
			t.Fatal(i.key, err)
		}

		rl, err := readTestItem(value)
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
	data []item,
	ic *IndexedCollection,
	tn tablename.TableName,
) {
	result, err := getListValues(ctx, data, ic, tn)
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
	data []item,
	ic *IndexedCollection,
	tn tablename.TableName,
) ([]TestItem, error) {
	result := make([]TestItem, 0, len(data))
	i := 0
	for r, err := range ic.List(ctx, &logtesting.MockLog{}, tn) {
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
