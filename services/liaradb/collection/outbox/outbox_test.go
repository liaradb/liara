package outbox

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
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/record"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/util/testing/logtesting"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestOutbox(t *testing.T) {
	storagetesting.SyncTest(t, 9, 256, testOutbox)
}

func testOutbox(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	o := New(fixed.New(record.New(s.Storage), btree.NewCursor(s.Storage)))
	tid := value.NewTenantID()

	data := createData()
	slices.Reverse(data)

	if err := insertData(ctx, &logtesting.MockLog{}, tid, o, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, tid, o, data)
	testList(ctx, t, tid, data, o)

	synctest.Wait()
}

func TestOutbox__LargeBuffer(t *testing.T) {
	storagetesting.SyncTest(t, 4, 256, testOutbox__LargeBuffer)
}

func testOutbox__LargeBuffer(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	o := New(fixed.New(record.New(s.Storage), btree.NewCursor(s.Storage)))
	tid := value.NewTenantID()

	data := createData()

	if err := insertData(ctx, &logtesting.MockLog{}, tid, o, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, tid, o, data)
	testList(ctx, t, tid, data, o)

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

func insertData(ctx context.Context, tx fixed.Transaction, tid value.TenantID, o *Outbox, data []item) error {
	for _, i := range data {
		if err := o.Insert(ctx, tx, tid, i.value.ID(), i.value); err != nil {
			return err
		}
	}
	return nil
}

func testGet(
	ctx context.Context,
	t *testing.T,
	tid value.TenantID,
	o *Outbox,
	data []item,
) {
	for _, i := range data {
		value, err := o.Get(ctx, &logtesting.MockLog{}, tid, i.value.ID())
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
	tid value.TenantID,
	data []item,
	o *Outbox,
) {
	result, err := getListValues(ctx, tid, data, o)
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
	tid value.TenantID,
	data []item,
	o *Outbox,
) ([]entity.Outbox, error) {
	result := make([]entity.Outbox, 0, len(data))
	i := 0
	for value, err := range o.List(ctx, &logtesting.MockLog{}, tid) {
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
