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
	// TODO: This is flaky on insert when buffer count is 5
	// s := storagetesting.CreateStorage(t, 5, 84)
	s := storagetesting.CreateStorage(t, 7, 296)
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

func TestTenant__LargeBuffer(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testTenant__LargeBuffer)
}

func testTenant__LargeBuffer(t *testing.T) {
	ctx := t.Context()
	s := storagetesting.CreateStorage(t, 2, 1024)
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

func createData() map[string]*entity.Tenant {
	return map[string]*entity.Tenant{
		"1": entity.NewTenant(value.NewTenantID(), value.TenantName{}),
		"2": entity.NewTenant(value.NewTenantID(), value.TenantName{}),
		"3": entity.NewTenant(value.NewTenantID(), value.TenantName{}),
		"4": entity.NewTenant(value.NewTenantID(), value.TenantName{}),
		"5": entity.NewTenant(value.NewTenantID(), value.TenantName{}),
		"6": entity.NewTenant(value.NewTenantID(), value.TenantName{}),
		"7": entity.NewTenant(value.NewTenantID(), value.TenantName{}),
		"8": entity.NewTenant(value.NewTenantID(), value.TenantName{}),
		"9": entity.NewTenant(value.NewTenantID(), value.TenantName{}),
	}
}

func insertData(ctx context.Context, o *Tenant, n tablename.TableName, data map[string]*entity.Tenant) error {
	for _, v := range data {
		if err := o.Set(ctx, n, v.ID(), v); err != nil {
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
	data map[string]*entity.Tenant,
) {
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

func testList(
	ctx context.Context,
	t *testing.T,
	data map[string]*entity.Tenant,
	o *Tenant,
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
	data map[string]*entity.Tenant,
	o *Tenant,
	n tablename.TableName,
	pid value.PartitionID,
) ([]entity.Tenant, error) {
	result := make([]entity.Tenant, 0, len(data))
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

func createSortedValues(data map[string]*entity.Tenant) []entity.Tenant {
	type tuple struct {
		key   key.Key
		value *entity.Tenant
	}

	tuples := make([]tuple, 0, len(data))
	for _, v := range data {
		tuples = append(tuples, tuple{key.NewKey(v.ID().Bytes()), v})
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
