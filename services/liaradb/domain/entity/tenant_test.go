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

func TestTenant_ReadWrite(t *testing.T) {
	tid := value.NewTenantID()
	name := value.NewTenantName("name")
	tn := NewTenant(tid, name)

	data := make([]byte, TenantSize+2)
	data0 := tn.Write(data)

	if l := len(data0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	tn1 := &Tenant{}
	data1 := tn1.Read(data)
	if l := len(data1); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	if *tn1 != *tn {
		t.Errorf("incorrect result: %v, expected: %v", *tn1, *tn)
	}
}
