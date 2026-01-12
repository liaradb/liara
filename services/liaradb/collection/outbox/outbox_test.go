package outbox

import (
	"context"
	"slices"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestOutbox(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testOutbox)
}

func testOutbox(t *testing.T) {
	ctx := t.Context()
	// TODO: This is flaky on insert when buffer count is 5
	// s := storagetesting.CreateStorage(t, 5, 84)
	s := storagetesting.CreateStorage(t, 6, 84)
	o := New(s, btree.NewCursor(s))
	n := tablename.New("testfile")

	data := createData()

	if err := insertData(ctx, o, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, o, n, data)
	testList(ctx, t, data, o, n)

	synctest.Wait()
}

func TestOutbox__LargeBuffer(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testOutbox__LargeBuffer)
}

func testOutbox__LargeBuffer(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 2, 256)
	o := New(s, btree.NewCursor(s))
	n := tablename.New("testfile")

	data := createData()

	if err := insertData(ctx, o, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, o, n, data)
	testList(ctx, t, data, o, n)

	synctest.Wait()
}

func createData() map[string][]byte {
	return map[string][]byte{
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

func insertData(ctx context.Context, o *Outbox, n tablename.TableName, data map[string][]byte) error {
	for k, v := range data {
		if err := o.Set(ctx, n, key.NewKey([]byte(k)), v); err != nil {
			return err
		}
	}
	return nil
}

func testGet(ctx context.Context, t *testing.T, kv *Outbox, n tablename.TableName, data map[string][]byte) {
	for k, v := range data {
		value, err := kv.Get(ctx, n, key.NewKey([]byte(k)))
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

func testList(ctx context.Context, t *testing.T, data map[string][]byte, o *Outbox, n tablename.TableName) {
	result, err := getListValues(ctx, data, o, n)
	if err != nil {
		t.Fatal(err)
	}

	want := createSortedValues(data)
	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func getListValues(ctx context.Context, data map[string][]byte, o *Outbox, n tablename.TableName) ([]string, error) {
	result := make([]string, 0, len(data))
	i := 0
	for value, err := range o.List(ctx, n) {
		if err != nil {
			return nil, err
		}

		result = append(result, string(value))
		i++
	}
	return result, nil
}

func createSortedValues(data map[string][]byte) []string {
	type tuple struct {
		key   key.Key
		value []byte
	}

	tuples := make([]tuple, 0, len(data))
	for k, v := range data {
		tuples = append(tuples, tuple{key.NewKey([]byte(k)), v})
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
