package keyvalue

import (
	"context"
	"slices"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestKeyValue(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testKeyValue)
}

func testKeyValue(t *testing.T) {
	ctx := t.Context()
	// TODO: This is flaky on insert when buffer count is 4
	// s := storagetesting.CreateStorage(t, 4, 64)
	s := storagetesting.CreateStorage(t, 5, 64)
	kv := New(s)
	n := tablename.New("testfile")

	data := createData()

	if err := insertData(ctx, kv, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, kv, n, data)
	testList(ctx, t, data, kv, n)

	synctest.Wait()
}

func TestKeyValue__LargeBuffer(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testKeyValue__LargeBuffer)
}

func testKeyValue__LargeBuffer(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 2, 256)
	kv := New(s)
	n := tablename.New("testfile")

	data := createData()

	if err := insertData(ctx, kv, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, kv, n, data)
	testList(ctx, t, data, kv, n)

	synctest.Wait()
}

func createData() map[value.Key][]byte {
	return map[value.Key][]byte{
		"1": []byte("a"),
		"2": []byte("b"),
		"3": []byte("c"),
		"4": []byte("d"),
		"5": []byte("e"),
		"6": []byte("f"),
		"7": []byte("g"),
		"8": []byte("h"),
		"9": []byte("i"),
	}
}

func insertData(ctx context.Context, kv *KeyValue, n tablename.TableName, data map[value.Key][]byte) error {
	for k, v := range data {
		if err := kv.Set(ctx, n, k, v); err != nil {
			return err
		}
	}
	return nil
}

func testGet(ctx context.Context, t *testing.T, kv *KeyValue, n tablename.TableName, data map[value.Key][]byte) {
	for k, v := range data {
		value, err := kv.Get(ctx, n, k)
		if err != nil {
			t.Fatal(k, err)
		}

		want := string(v)
		result := string(value)
		if result != want {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	}
}

func testList(ctx context.Context, t *testing.T, data map[value.Key][]byte, kv *KeyValue, n tablename.TableName) {
	result, err := getListValues(ctx, data, kv, n)
	if err != nil {
		t.Fatal(err)
	}

	want := createSortedValues(data)
	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func getListValues(ctx context.Context, data map[value.Key][]byte, kv *KeyValue, n tablename.TableName) ([]string, error) {
	result := make([]string, 0, len(data))
	i := 0
	for value, err := range kv.List(ctx, n) {
		if err != nil {
			return nil, err
		}

		result = append(result, string(value))
		i++
	}
	return result, nil
}

func createSortedValues(data map[value.Key][]byte) []string {
	type tuple struct {
		key   value.Key
		value []byte
	}

	tuples := make([]tuple, 0, len(data))
	for k, v := range data {
		tuples = append(tuples, tuple{k, v})
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
