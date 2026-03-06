package tenant

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

func TestTenant(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTenant)
}

func testTenant(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 5, 296)
	o := New(s, btree.NewCursor(s))
	n := tablename.NewFromString("testfile")

	data := createData()
	slices.Reverse(data)

	if err := insertData(ctx, o, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, o, n, data)
	testList(ctx, t, data, o, n)

	synctest.Wait()
}

func TestTenant__LargeBuffer(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTenant__LargeBuffer)
}

func testTenant__LargeBuffer(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 2, 1024)
	o := New(s, btree.NewCursor(s))
	n := tablename.NewFromString("testfile")

	data := createData()

	if err := insertData(ctx, o, n, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, o, n, data)
	testList(ctx, t, data, o, n)

	synctest.Wait()
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

func insertData(ctx context.Context, o *Tenant, n tablename.TableName, data []item) error {
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
	kv *Tenant,
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
	o *Tenant,
	n tablename.TableName,
) {
	result, err := getListValues(ctx, data, o, n)
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
	o *Tenant,
	n tablename.TableName,
) ([]entity.Tenant, error) {
	result := make([]entity.Tenant, 0, len(data))
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
