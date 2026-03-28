package keyvalue

import (
	"context"
	"slices"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestKeyValue(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip bool
		data []item
	}{
		"should insert data": {
			data: createData(),
		},
		"should insert data in reverse": {
			data: createDataReverse(),
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			synctest.Test(t, func(t *testing.T) {
				ctx := t.Context()

				s := storagetesting.CreateStorage(t, 8, 84)
				kv := New(s, btree.NewCursor(s))
				tn := tablename.NewFromString("testfile")
				pid := value.NewPartitionID(0)

				if err := insertData(ctx, kv, tn, pid, c.data); err != nil {
					t.Fatal(err)
				}

				testGet(ctx, t, kv, tn, pid, c.data)
				testList(ctx, t, c.data, kv, tn, pid)

				synctest.Wait()
			})
		})
	}
}

func TestKeyValue__LargeBuffer(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testKeyValue__LargeBuffer)
}

func testKeyValue__LargeBuffer(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 2, 256)
	kv := New(s, btree.NewCursor(s))
	tn := tablename.NewFromString("testfile")
	pid := value.NewPartitionID(0)

	data := createData()

	if err := insertData(ctx, kv, tn, pid, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, kv, tn, pid, data)
	testList(ctx, t, data, kv, tn, pid)

	synctest.Wait()
}

type item struct {
	key   string
	value []byte
}

func createDataReverse() []item {
	data := createData()
	slices.Reverse(data)
	return data
}

func createData() []item {
	return []item{
		{"1", []byte("a")},
		{"2", []byte("b")},
		{"3", []byte("c")},
		{"4", []byte("d")},
		{"5", []byte("e")},
		{"6", []byte("f")},
		{"7", []byte("g")},
		{"8", []byte("h")},
		{"9", []byte("i")},
	}
}

func insertData(ctx context.Context, kv *KeyValue, tn tablename.TableName, pid value.PartitionID, data []item) error {
	for _, i := range data {
		if err := kv.Set(ctx, tn, pid, key.NewKey([]byte(i.key)), i.value); err != nil {
			return err
		}
		synctest.Wait()
	}
	return nil
}

func testGet(ctx context.Context, t *testing.T, kv *KeyValue, tn tablename.TableName, pid value.PartitionID, data []item) {
	for _, i := range data {
		value, err := kv.Get(ctx, tn, pid, key.NewKey([]byte(i.key)))
		if err != nil {
			t.Fatal(i.key, err)
		}

		want := string(i.value)
		result := string(value)
		if result != want {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	}
}

func testList(ctx context.Context, t *testing.T, data []item, kv *KeyValue, tn tablename.TableName, pid value.PartitionID) {
	result, err := getListValues(ctx, data, kv, tn, pid)
	if err != nil {
		t.Fatal(err)
	}

	want := createSortedValues(data)
	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func getListValues(ctx context.Context, data []item, kv *KeyValue, tn tablename.TableName, pid value.PartitionID) ([]string, error) {
	result := make([]string, 0, len(data))
	i := 0
	for value, err := range kv.List(ctx, tn, pid) {
		if err != nil {
			return nil, err
		}

		result = append(result, string(value))
		i++
	}
	return result, nil
}

func createSortedValues(data []item) []string {
	type tuple struct {
		key   key.Key
		value []byte
	}

	tuples := make([]tuple, 0, len(data))
	for _, i := range data {
		tuples = append(tuples, tuple{key.NewKey([]byte(i.key)), i.value})
	}
	slices.SortFunc(tuples, func(a, b tuple) int {
		return strings.Compare(a.key.String(), b.key.String())
	})
	want := make([]string, 0, len(data))
	for _, t := range tuples {
		want = append(want, string(t.value))
	}
	return want
}
