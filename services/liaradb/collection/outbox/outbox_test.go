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
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
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
	s := storagetesting.CreateStorage(t, 7, 110)
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

func createData() map[string]*entity.Outbox {
	return map[string]*entity.Outbox{
		"1": entity.NewOutbox(value.NewOutboxID(), value.NewPartitionRange()),
		"2": entity.NewOutbox(value.NewOutboxID(), value.NewPartitionRange()),
		"3": entity.NewOutbox(value.NewOutboxID(), value.NewPartitionRange()),
		"4": entity.NewOutbox(value.NewOutboxID(), value.NewPartitionRange()),
		"5": entity.NewOutbox(value.NewOutboxID(), value.NewPartitionRange()),
		"6": entity.NewOutbox(value.NewOutboxID(), value.NewPartitionRange()),
		"7": entity.NewOutbox(value.NewOutboxID(), value.NewPartitionRange()),
		"8": entity.NewOutbox(value.NewOutboxID(), value.NewPartitionRange()),
		"9": entity.NewOutbox(value.NewOutboxID(), value.NewPartitionRange()),
	}
}

func insertData(ctx context.Context, o *Outbox, n tablename.TableName, data map[string]*entity.Outbox) error {
	for _, v := range data {
		if err := o.Set(ctx, n, v.ID(), v); err != nil {
			return err
		}
	}
	return nil
}

func testGet(ctx context.Context, t *testing.T, kv *Outbox, n tablename.TableName, data map[string]*entity.Outbox) {
	for k, v := range data {
		value, err := kv.Get(ctx, n, v.ID())
		if err != nil {
			t.Fatal(k, err)
		}

		if *value != *v {
			t.Errorf("incorrect result: %v, expected: %v", *value, *v)
		}
	}
}

func testList(ctx context.Context, t *testing.T, data map[string]*entity.Outbox, o *Outbox, n tablename.TableName) {
	result, err := getListValues(ctx, data, o, n)
	if err != nil {
		t.Fatal(err)
	}

	want := createSortedValues(data)
	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func getListValues(ctx context.Context, data map[string]*entity.Outbox, o *Outbox, n tablename.TableName) ([]entity.Outbox, error) {
	result := make([]entity.Outbox, 0, len(data))
	i := 0
	for value, err := range o.List(ctx, n) {
		if err != nil {
			return nil, err
		}

		result = append(result, *value)
		i++
	}
	return result, nil
}

func createSortedValues(data map[string]*entity.Outbox) []entity.Outbox {
	type tuple struct {
		key   key.Key
		value *entity.Outbox
	}

	tuples := make([]tuple, 0, len(data))
	for _, v := range data {
		tuples = append(tuples, tuple{key.NewKey(v.ID().Bytes()), v})
	}
	slices.SortFunc(tuples, func(a, b tuple) int {
		return strings.Compare(a.key.String(), b.key.String())
	})
	want := make([]entity.Outbox, 0, len(data))
	for _, t := range tuples {
		want = append(want, *t.value)
	}
	return want
}
