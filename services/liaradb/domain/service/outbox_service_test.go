package service

import (
	"testing"
	"testing/synctest"

	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/heap"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/infrastructure/idempotency"
	"github.com/liaradb/liaradb/domain/infrastructure/outbox"
	"github.com/liaradb/liaradb/domain/infrastructure/tenant"
	"github.com/liaradb/liaradb/domain/value"
)

func TestOutboxService_CreateOutbox(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		m, s := createManager(t, 1024, 1024)
		c := btree.NewCursor(s)
		ic := indexed.New(heap.New(s), c)
		tc := tenant.New(ic)
		idc := idempotency.New(ic, c)
		oc := outbox.New(ic)

		ts := NewTenantService(m, tc)
		os := NewOutboxService(m, tc, idc, oc)

		tid, err := ts.Create(t.Context(), command.CreateTenant{
			TenantName: value.NewTenantName("name"),
		})
		if err != nil {
			t.Fatal(err)
		}

		pid := value.NewPartitionID(2)
		oid, err := os.CreateOutbox(t.Context(), tid, pid)
		if err != nil {
			t.Fatal(err)
		}

		o, err := os.GetOutbox(t.Context(), tid, oid)
		if err != nil {
			t.Fatal(err)
		}

		if i := o.ID(); i != oid {
			t.Errorf("incorrect id: %v, expected: %v", i, oid)
		}
		if i := o.PartitionID(); i != pid {
			t.Errorf("incorrect partition id: %v, expected: %v", i, pid)
		}
		if v, want := o.GlobalVersion(), value.NewGlobalVersion(0); v != want {
			t.Errorf("incorrect global version: %v, expected: %v", v, want)
		}
		if v, want := o.SubVersion(), value.NewSubVersion(0); v != want {
			t.Errorf("incorrect sub version: %v, expected: %v", v, want)
		}

		synctest.Wait()
	})
}

func TestOutboxService_Update(t *testing.T) {
	t.Parallel()

	synctest.Test(t, func(t *testing.T) {
		m, s := createManager(t, 1024, 1024)
		c := btree.NewCursor(s)
		ic := indexed.New(heap.New(s), c)
		tc := tenant.New(ic)
		idc := idempotency.New(ic, c)
		oc := outbox.New(ic)

		ts := NewTenantService(m, tc)
		os := NewOutboxService(m, tc, idc, oc)

		tid, err := ts.Create(t.Context(), command.CreateTenant{
			TenantName: value.NewTenantName("name"),
		})
		if err != nil {
			t.Fatal(err)
		}

		pid := value.NewPartitionID(2)
		oid, err := os.CreateOutbox(t.Context(), tid, pid)
		if err != nil {
			t.Fatal(err)
		}

		gv := value.NewGlobalVersion(10)
		sv := value.NewSubVersion(200)
		if err := os.UpdateOutboxPosition(t.Context(), tid, oid, gv, sv); err != nil {
			t.Fatal(err)
		}

		o, err := os.GetOutbox(t.Context(), tid, oid)
		if err != nil {
			t.Fatal(err)
		}

		if i := o.ID(); i != oid {
			t.Errorf("incorrect id: %v, expected: %v", i, oid)
		}
		if i := o.PartitionID(); i != pid {
			t.Errorf("incorrect partition id: %v, expected: %v", i, pid)
		}
		if v := o.GlobalVersion(); v != gv {
			t.Errorf("incorrect global version: %v, expected: %v", v, gv)
		}
		if v := o.SubVersion(); v != sv {
			t.Errorf("incorrect sub version: %v, expected: %v", v, sv)
		}

		synctest.Wait()
	})
}
