package idempotency

import (
	"context"
	"slices"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage/storagetesting"
)

func TestIdempotency(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testIdempotency)
}

func testIdempotency(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 6, 110)
	o := New(s, btree.NewCursor(s))
	n := tablename.NewFromString("testfile")
	pid := value.NewPartitionID(0)

	data := createData()
	slices.Reverse(data)

	if err := insertData(ctx, o, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, o, n, pid, data)
	testList(ctx, t, data, o, n, pid)

	synctest.Wait()
}

func TestRequestLog__LargeBuffer(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testRequestLog__LargeBuffer)
}

func testRequestLog__LargeBuffer(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 2, 256)
	o := New(s, btree.NewCursor(s))
	n := tablename.NewFromString("testfile")
	pid := value.NewPartitionID(0)

	data := createData()

	if err := insertData(ctx, o, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, o, n, pid, data)
	testList(ctx, t, data, o, n, pid)

	synctest.Wait()
}

type item struct {
	key   string
	value *entity.RequestLog
}

func createData() []item {
	return []item{
		{"1", entity.NewRequestLog(value.NewRequestID(), time.Now().Add(0*time.Second))},
		{"2", entity.NewRequestLog(value.NewRequestID(), time.Now().Add(1*time.Second))},
		{"3", entity.NewRequestLog(value.NewRequestID(), time.Now().Add(2*time.Second))},
		{"4", entity.NewRequestLog(value.NewRequestID(), time.Now().Add(3*time.Second))},
		{"5", entity.NewRequestLog(value.NewRequestID(), time.Now().Add(4*time.Second))},
		{"6", entity.NewRequestLog(value.NewRequestID(), time.Now().Add(5*time.Second))},
		{"7", entity.NewRequestLog(value.NewRequestID(), time.Now().Add(6*time.Second))},
		{"8", entity.NewRequestLog(value.NewRequestID(), time.Now().Add(7*time.Second))},
		{"9", entity.NewRequestLog(value.NewRequestID(), time.Now().Add(8*time.Second))},
	}
}

func insertData(ctx context.Context, o *Idempotency, n tablename.TableName, data []item) error {
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
	kv *Idempotency,
	n tablename.TableName,
	pid value.PartitionID,
	data []item,
) {
	for _, i := range data {
		value, err := kv.Get(ctx, n, pid, i.value.ID())
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
	o *Idempotency,
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
	o *Idempotency,
	n tablename.TableName,
	pid value.PartitionID,
) ([]entity.RequestLog, error) {
	result := make([]entity.RequestLog, 0, len(data))
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

func createSortedValues(data []item) []entity.RequestLog {
	type tuple struct {
		key   key.Key
		value *entity.RequestLog
	}

	tuples := make([]tuple, 0, len(data))
	for _, i := range data {
		tuples = append(tuples, tuple{key.NewKey(i.value.ID().Bytes()), i.value})
	}
	slices.SortFunc(tuples, func(a, b tuple) int {
		return strings.Compare(a.key.String(), b.key.String())
	})
	want := make([]entity.RequestLog, 0, len(data))
	for _, t := range tuples {
		want = append(want, *t.value)
	}
	return want
}
