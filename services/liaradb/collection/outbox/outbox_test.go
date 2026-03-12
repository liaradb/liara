package outbox

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/google/uuid"
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestOutbox(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testOutbox)
}

func testOutbox(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 7, 110)
	o := New(s, btree.NewCursor(s))
	n := tablename.NewFromString("testfile")
	pid := value.NewPartitionID(0)

	data := createData()
	slices.Reverse(data)

	if err := insertData(ctx, o, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, o, n, data)
	testList(ctx, t, data, o, n, pid)

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
	n := tablename.NewFromString("testfile")
	pid := value.NewPartitionID(0)

	data := createData()

	if err := insertData(ctx, o, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, o, n, data)
	testList(ctx, t, data, o, n, pid)

	synctest.Wait()
}

type item struct {
	key   string
	value *entity.Outbox
}

func createData() []item {
	count := 9
	data := uuid.UUID{}
	items := make([]item, 0, count)
	for i := range count {
		data[15] = byte(i) + 1
		rid, _ := value.NewOutboxIDFromString(data.String())
		items = append(items, item{fmt.Sprintf("%v", i+1), entity.NewOutbox(rid, value.NewPartitionRange(value.NewPartitionID(0), value.NewPartitionID(0)))})
	}
	return items
}

func insertData(ctx context.Context, o *Outbox, n tablename.TableName, data []item) error {
	for _, i := range data {
		if err := o.Set(ctx, n, i.value.ID(), i.value); err != nil {
			return err
		}
	}
	return nil
}

func testGet(
	ctx context.Context,
	t *testing.T,
	kv *Outbox,
	n tablename.TableName,
	data []item,
) {
	for _, i := range data {
		value, err := kv.Get(ctx, n, i.value.ID())
		if err != nil {
			t.Fatal(i.key, err)
		}

		if *value != *i.value {
			t.Errorf("incorrect result: %v, expected: %v", *value, *i.value)
		}
	}
}

func testList(
	ctx context.Context,
	t *testing.T,
	data []item,
	o *Outbox,
	n tablename.TableName,
	pid value.PartitionID,
) {
	result, err := getListValues(ctx, data, o, n, pid)
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
	o *Outbox,
	n tablename.TableName,
	pid value.PartitionID,
) ([]entity.Outbox, error) {
	result := make([]entity.Outbox, 0, len(data))
	i := 0
	for value, err := range o.List(ctx, n, pid) {
		if err != nil {
			return nil, err
		}

		result = append(result, *value)
		i++
	}
	return result, nil
}

func createSortedValues(data []item) []entity.Outbox {
	type tuple struct {
		key   key.Key
		value *entity.Outbox
	}

	tuples := make([]tuple, 0, len(data))
	for _, i := range data {
		tuples = append(tuples, tuple{key.NewKey(i.value.ID().Bytes()), i.value})
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
