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
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/record"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/util/testing/logtesting"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestIdempotency(t *testing.T) {
	storagetesting.SyncTest(t, 8, 256, testIdempotency)
}

func testIdempotency(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	idp := New(fixed.New(record.New(s.Storage), btree.NewCursor(s.Storage)))
	n := tablename.NewTenantFolder(value.NewTenantID()).RequestLog()

	data := createData()
	slices.Reverse(data)

	if err := insertData(ctx, &logtesting.MockLog{}, idp, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, idp, n, data)
	testList(ctx, t, data, idp, n)

	synctest.Wait()
}

func TestRequestLog__LargeBuffer(t *testing.T) {
	storagetesting.SyncTest(t, 4, 256, testRequestLog__LargeBuffer)
}

func testRequestLog__LargeBuffer(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	idp := New(fixed.New(record.New(s.Storage), btree.NewCursor(s.Storage)))
	n := tablename.NewTenantFolder(value.NewTenantID()).RequestLog()

	data := createData()

	if err := insertData(ctx, &logtesting.MockLog{}, idp, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, idp, n, data)
	testList(ctx, t, data, idp, n)

	synctest.Wait()
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

func insertData(ctx context.Context, tx fixed.Transaction, o *Idempotency, tn tablename.TableName, data []item) error {
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
	idp *Idempotency,
	n tablename.TableName,
	data []item,
) {
	for _, i := range data {
		value, err := idp.Get(ctx, &logtesting.MockLog{}, n, i.value.ID())
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
	idp *Idempotency,
	n tablename.TableName,
) {
	result, err := getListValues(ctx, data, idp, n)
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
) ([]entity.RequestLog, error) {
	result := make([]entity.RequestLog, 0, len(data))
	i := 0
	for value, err := range o.List(ctx, &logtesting.MockLog{}, n) {
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
