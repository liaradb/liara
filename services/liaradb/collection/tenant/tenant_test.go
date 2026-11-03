package tenant

import (
	"context"
	"slices"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/transaction/log"
	"github.com/liaradb/liaradb/util/testing/logtesting"
	"github.com/liaradb/liaradb/util/testing/storagetesting"
)

func TestTenant(t *testing.T) {
	// TODO: Why can't this page size be smaller?
	storagetesting.SyncTest(t, 5, 302, testTenant)
}

func testTenant(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	l := log.New(4096, 100, 4096, 100, s.FSys, "dir")
	if err := l.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	o := New(s.Storage, btree.NewCursor(s.Storage), l)

	data := createData()
	slices.Reverse(data)

	if err := insertData(ctx, &logtesting.MockLog{}, o, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, o, data)
	testList(ctx, t, data, o)

	synctest.Wait()
}

func TestTenant__LargeBuffer(t *testing.T) {
	storagetesting.SyncTest(t, 3, 1024, testTenant__LargeBuffer)
}

func testTenant__LargeBuffer(t *testing.T, s storagetesting.Storage) {
	ctx := t.Context()
	l := log.New(4096, 100, 4096, 100, s.FSys, "dir")
	if err := l.Run(t.Context()); err != nil {
		t.Fatal(err)
	}

	if err := l.StartWriter(); err != nil {
		t.Fatal(err)
	}

	o := New(s.Storage, btree.NewCursor(s.Storage), l)

	data := createData()

	if err := insertData(ctx, &logtesting.MockLog{}, o, data); err != nil {
		t.Fatal(err)
	}

	testGet(ctx, t, o, data)
	testList(ctx, t, data, o)

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

func insertData(ctx context.Context, l span.Log, o *Tenant, data []item) error {
	for _, i := range data {
		if err := o.Insert(ctx, l, i.value.ID(), i.value); err != nil {
			return err
		}
	}
	return nil
}

func testGet(
	ctx context.Context,
	t *testing.T,
	kv *Tenant,
	data []item,
) {
	for _, i := range data {
		value, err := kv.Get(ctx, &logtesting.MockLog{}, i.value.ID())
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
) {
	result, err := getListValues(ctx, data, o)
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
) ([]entity.Tenant, error) {
	result := make([]entity.Tenant, 0, len(data))
	i := 0
	for value, err := range o.List(ctx, &logtesting.MockLog{}) {
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
