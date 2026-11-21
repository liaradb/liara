package idempotency

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"testing/synctest"
	"time"
	"uuid"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/heap"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/infrastructure/tenantfolder"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
	"github.com/liaradb/liaradb/util/testing/txtesting"
)

func TestIdempotency(t *testing.T) {
	storagetesting.SyncTest(t, 8, 256, testIdempotency)
}

func testIdempotency(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	c := btree.NewCursor(s.Storage)
	idp := New(indexed.New(heap.New(), c), c)
	n := tenantfolder.NewTenantFolder(value.NewTenantID()).RequestLog()

	data := createData()
	slices.Reverse(data)

	tx := txtesting.New(s.Storage)
	defer tx.Release()

	if err := insertData(ctx, tx, idp, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, tx, idp, n, data)
	testList(ctx, t, tx, data, idp, n)

	synctest.Wait()
}

func TestRequestLog__LargeBuffer(t *testing.T) {
	storagetesting.SyncTest(t, 4, 256, testRequestLog__LargeBuffer)
}

func testRequestLog__LargeBuffer(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	c := btree.NewCursor(s.Storage)
	idp := New(indexed.New(heap.New(), c), c)
	n := tenantfolder.NewTenantFolder(value.NewTenantID()).RequestLog()

	data := createData()

	tx := txtesting.New(s.Storage)
	defer tx.Release()

	if err := insertData(ctx, tx, idp, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, tx, idp, n, data)
	testList(ctx, t, tx, data, idp, n)

	synctest.Wait()
}

func TestIdempotency_Test(t *testing.T) {
	storagetesting.SyncTest(t, 8, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		c := btree.NewCursor(s.Storage)
		idp := New(indexed.New(heap.New(), c), c)
		tn := tenantfolder.NewTenantFolder(value.NewTenantID()).RequestLog()

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		if ok, err := idp.Test(ctx, tx, tn, value.NewRequestID()); err != nil {
			t.Fatal(err)
		} else if !ok {
			t.Error("should be true")
		}

		rqid := value.NewRequestID()
		rl := entity.NewRequestLog(rqid, value.NewTime(time.Now()))

		if err := idp.Insert(ctx, tx, tn, rqid, rl); err != nil {
			t.Fatal(err)
		}

		if ok, err := idp.Test(ctx, tx, tn, rqid); err != nil {
			t.Fatal(err)
		} else if ok {
			t.Error("should be false")
		}
	})
}

func TestIdempotency_Get__Error(t *testing.T) {
	storagetesting.SyncTest(t, 4, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		c := btree.NewCursor(s.Storage)
		idp := New(indexed.New(heap.New(), c), c)
		tid := value.NewTenantID()
		tn := tenantfolder.NewTenantFolder(tid).Outbox()
		rqid := value.NewRequestID()

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		if _, err := idp.Get(ctx, tx, tn, rqid); err == nil {
			t.Fatal("should return error")
		}
	})
}

type item struct {
	key   string
	value *entity.RequestLog
}

func createData() []item {
	count := 9
	data := uuid.UUID{}
	items := make([]item, 0, count)
	for i := range count {
		data[15] = byte(i) + 1
		rid, _ := value.NewRequestIDFromString(data.String())
		items = append(items, item{fmt.Sprintf("%v", i+1), entity.NewRequestLog(rid, value.NewTime(time.Now().Add(time.Duration(i)*time.Second)))})
	}
	return items
}

func insertData(ctx context.Context, tx indexed.Transaction, o *Idempotency, tn tablename.TableName, data []item) error {
	for _, i := range data {
		if err := o.Insert(ctx, tx, tn, i.value.ID(), i.value); err != nil {
			return err
		}
	}
	return nil
}

func testGet(
	ctx context.Context,
	t *testing.T,
	tx indexed.Transaction,
	idp *Idempotency,
	n tablename.TableName,
	data []item,
) {
	for _, i := range data {
		value, err := idp.Get(ctx, tx, n, i.value.ID())
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
	tx indexed.Transaction,
	data []item,
	idp *Idempotency,
	tn tablename.TableName,
) {
	result, err := getListValues(ctx, tx, data, idp, tn)
	if err != nil {
		t.Fatal(err)
	}

	want := createSortedValues(data)
	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}

	for range idp.List(ctx, tx, tn) {
		break
	}
}

func getListValues(
	ctx context.Context,
	tx indexed.Transaction,
	data []item,
	o *Idempotency,
	tn tablename.TableName,
) ([]entity.RequestLog, error) {
	result := make([]entity.RequestLog, 0, len(data))
	i := 0
	for value, err := range o.List(ctx, tx, tn) {
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
