package tenant

import (
	"context"
	"slices"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/heap"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
	"github.com/liaradb/liaradb/util/testing/txtesting"
)

func TestTenant(t *testing.T) {
	storagetesting.SyncTest(t, 5, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		tn := New(indexed.New(heap.New(), btree.NewCursor(s.Storage)))

		data := createData()
		slices.Reverse(data)

		tx := txtesting.New(s.Storage)

		if err := insertData(ctx, tx, tn, data); err != nil {
			t.Fatal(err)
		}

		testGet(ctx, t, tx, tn, data)
		testList(ctx, t, tx, data, tn)

		tx.Release()
	})
}

func TestTenant__LargeBuffer(t *testing.T) {
	storagetesting.SyncTest(t, 3, 1024, testTenant__LargeBuffer)
}

func testTenant__LargeBuffer(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	tn := New(indexed.New(heap.New(), btree.NewCursor(s.Storage)))

	data := createData()

	tx := txtesting.New(s.Storage)
	defer tx.Release()

	if err := insertData(ctx, tx, tn, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, tx, tn, data)
	testList(ctx, t, tx, data, tn)

	synctest.Wait()
}

func TestTenant_Get__Error(t *testing.T) {
	storagetesting.SyncTest(t, 4, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		tn := New(indexed.New(heap.New(), btree.NewCursor(s.Storage)))
		tx := txtesting.New(s.Storage)
		defer tx.Release()

		tid := value.NewTenantID()

		if _, err := tn.Get(ctx, tx, tid); err == nil {
			t.Fatal("should return error")
		}
	})
}

func TestTenant_Replace(t *testing.T) {
	storagetesting.SyncTest(t, 10, 256, func(t *testing.T, s storagetesting.Storage) {
		ctx := t.Context()
		tn := New(indexed.New(heap.New(), btree.NewCursor(s.Storage)))

		tx := txtesting.New(s.Storage)
		defer tx.Release()

		id := value.NewTenantID()
		t0 := entity.NewTenant(id, value.NewTenantName("Tenant 1"))

		if err := tn.Insert(ctx, tx, t0); err != nil {
			t.Fatal(err)
		}

		t1, err := tn.Get(ctx, tx, id)
		if err != nil {
			t.Fatal(err)
		}

		if *t1 != *t0 {
			t.Errorf("incorrect result: %v, expected: %v", *t1, *t0)
		}

		t1.Rename(value.NewTenantName("Tenant 2"))

		if err := tn.Replace(ctx, tx, t1); err != nil {
			t.Fatal(err)
		}

		t2, err := tn.Get(ctx, tx, id)
		if err != nil {
			t.Fatal(err)
		}

		if *t2 != *t1 {
			t.Errorf("incorrect result: %v, expected: %v", *t2, *t1)
		}
	})
}

type item struct {
	key   string
	value *entity.Tenant
}

func createData() []item {
	return []item{
		{"1", entity.NewTenant(value.NewTenantID(), value.TenantName{})},
		{"2", entity.NewTenant(value.NewTenantID(), value.TenantName{})},
		{"3", entity.NewTenant(value.NewTenantID(), value.TenantName{})},
		{"4", entity.NewTenant(value.NewTenantID(), value.TenantName{})},
		{"5", entity.NewTenant(value.NewTenantID(), value.TenantName{})},
		{"6", entity.NewTenant(value.NewTenantID(), value.TenantName{})},
		{"7", entity.NewTenant(value.NewTenantID(), value.TenantName{})},
		{"8", entity.NewTenant(value.NewTenantID(), value.TenantName{})},
		{"9", entity.NewTenant(value.NewTenantID(), value.TenantName{})},
	}
}

func insertData(ctx context.Context, tx indexed.Transaction, o *Tenant, data []item) error {
	for _, i := range data {
		if err := o.Insert(ctx, tx, i.value); err != nil {
			return err
		}
	}
	return nil
}

func testGet(
	ctx context.Context,
	t *testing.T,
	tx indexed.Transaction,
	tn *Tenant,
	data []item,
) {
	for _, i := range data {
		value, err := tn.Get(ctx, tx, i.value.ID())
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
	tn *Tenant,
) {
	result, err := getListValues(ctx, tx, data, tn)
	if err != nil {
		t.Fatal(err)
	}

	want := createSortedValues(data)
	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}

	for range tn.List(ctx, tx) {
		break
	}
}

func getListValues(
	ctx context.Context,
	tx indexed.Transaction,
	data []item,
	tn *Tenant,
) ([]entity.Tenant, error) {
	result := make([]entity.Tenant, 0, len(data))
	i := 0
	for value, err := range tn.List(ctx, tx) {
		if err != nil {
			return nil, err
		}

		result = append(result, *value)
		i++
	}
	return result, nil
}

func createSortedValues(data []item) []entity.Tenant {
	type tuple struct {
		key   key.Key
		value *entity.Tenant
	}

	tuples := make([]tuple, 0, len(data))
	for _, i := range data {
		tuples = append(tuples, tuple{key.NewKey(i.value.ID().Bytes()), i.value})
	}
	slices.SortFunc(tuples, func(a, b tuple) int {
		return strings.Compare(a.key.String(), b.key.String())
	})
	want := make([]entity.Tenant, 0, len(data))
	for _, t := range tuples {
		want = append(want, *t.value)
	}
	return want
}
