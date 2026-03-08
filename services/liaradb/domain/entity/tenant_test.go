package entity

import (
	"testing"

	"github.com/liaradb/liaradb/domain/value"
)

func TestTenant(t *testing.T) {
	tid := value.NewTenantID()
	name := value.NewTenantName("name")
	tn := NewTenant(tid, name)

	if i := tn.ID(); i != tid {
		t.Errorf("incorrect id: %v, expected: %v", i, tid)
	}

	if v := tn.Version().Value(); v != 0 {
		t.Errorf("incorrect version: %v, expected: %v", v, 0)
	}

	if n := tn.Name(); n != name {
		t.Errorf("incorrect name: %v, expected: %v", n, name)
	}
}

func TestTenant_Rename(t *testing.T) {
	tid := value.NewTenantID()
	name := value.NewTenantName("name")
	tn := NewTenant(tid, name)

	name2 := value.NewTenantName("new name")
	tn.Rename(name2)

	if n := tn.Name(); n != name2 {
		t.Errorf("incorrect name: %v, expected: %v", n, name2)
	}

	if v := tn.Version().Value(); v != 1 {
		t.Errorf("incorrect version: %v, expected: %v", v, 1)
	}

}
