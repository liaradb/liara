package service

import (
	"slices"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/heap"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/infrastructure/tenant"
	"github.com/liaradb/liaradb/domain/value"
)

func TestTenantService_Create(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		m, _ := createManager(t, 1024, 1024)
		ic := indexed.New(heap.New())
		ts := NewTenantService(m, tenant.New(ic))

		name := value.NewTenantName("name")
		tid, err := ts.Create(t.Context(), command.CreateTenant{
			TenantName: name,
		})
		if err != nil {
			t.Fatal(err)
		}

		tn, err := ts.Get(t.Context(), tid)
		if err != nil {
			t.Fatal(err)
		}

		if i := tn.ID(); i != tid {
			t.Errorf("incorrect id: %v, expected: %v", i, i)
		}
		if n := tn.Name(); n != name {
			t.Errorf("incorrect name: %v, expected: %v", n, name)
		}

		synctest.Wait()
	})
}

func TestTenantService_Rename(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		m, _ := createManager(t, 1024, 1024)
		ic := indexed.New(heap.New())
		ts := NewTenantService(m, tenant.New(ic))

		tid, err := ts.Create(t.Context(), command.CreateTenant{
			TenantName: value.NewTenantName("name"),
		})
		if err != nil {
			t.Fatal(err)
		}

		name := value.NewTenantName("name 1")
		if err := ts.Rename(t.Context(), command.RenameTenant{
			TenantID:   tid,
			TenantName: name,
		}); err != nil {
			t.Fatal(err)
		}

		tn, err := ts.Get(t.Context(), tid)
		if err != nil {
			t.Fatal(err)
		}

		if i := tn.ID(); i != tid {
			t.Errorf("incorrect id: %v, expected: %v", i, i)
		}
		if n := tn.Name(); n != name {
			t.Errorf("incorrect name: %v, expected: %v", n, name)
		}

		synctest.Wait()
	})
}

func TestTenantService_List(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		m, _ := createManager(t, 1024, 1024)
		ic := indexed.New(heap.New())
		ts := NewTenantService(m, tenant.New(ic))

		tid0, err := ts.Create(t.Context(), command.CreateTenant{
			TenantName: value.NewTenantName("name 0"),
		})
		if err != nil {
			t.Fatal(err)
		}

		tid1, err := ts.Create(t.Context(), command.CreateTenant{
			TenantName: value.NewTenantName("name 1"),
		})
		if err != nil {
			t.Fatal(err)
		}

		tid2, err := ts.Create(t.Context(), command.CreateTenant{
			TenantName: value.NewTenantName("name 2"),
		})
		if err != nil {
			t.Fatal(err)
		}

		want := []entity.Tenant{
			*entity.NewTenant(tid0, value.NewTenantName("name 0")),
			*entity.NewTenant(tid1, value.NewTenantName("name 1")),
			*entity.NewTenant(tid2, value.NewTenantName("name 2"))}
		slices.SortFunc(want, func(a, b entity.Tenant) int {
			return strings.Compare(a.ID().String(), b.ID().String())
		})

		result := make([]entity.Tenant, 0, len(want))
		for tn, err := range ts.List(t.Context(), 0, 0) {
			if err != nil {
				t.Fatal(err)
			}

			result = append(result, *tn)
		}
		slices.SortFunc(result, func(a, b entity.Tenant) int {
			return strings.Compare(a.ID().String(), b.ID().String())
		})

		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}

		for range ts.List(t.Context(), 0, 0) {
			break
		}

		synctest.Wait()
	})
}
