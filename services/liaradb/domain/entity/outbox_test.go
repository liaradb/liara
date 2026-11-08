package entity

import (
	"testing"

	"github.com/liaradb/liaradb/domain/value"
)

func TestOutbox(t *testing.T) {
	t.Parallel()

	oid := value.NewOutboxID()
	pid := value.NewPartitionID(2)

	o := NewOutbox(oid, pid)

	if i := o.ID(); i != oid {
		t.Errorf("incorrect id: %v, expected: %v", i, oid)
	}

	if p := o.PartitionID(); p != pid {
		t.Errorf("incorrect partition id: %v, expected: %v", p, pid)
	}

	if v := o.GlobalVersion().Value(); v != 0 {
		t.Errorf("incorrect global version: %v, expected: %v", v, 0)
	}
}

func TestOutbox_UpdateGlobalVersion(t *testing.T) {
	t.Parallel()

	oid := value.NewOutboxID()
	pid := value.NewPartitionID(2)

	o := NewOutbox(oid, pid)

	if v := o.GlobalVersion().Value(); v != 0 {
		t.Errorf("incorrect global version: %v, expected: %v", v, 0)
	}

	gv := value.NewGlobalVersion(12345)

	o.UpdateGlobalVersion(gv)

	if v := o.GlobalVersion(); v != gv {
		t.Errorf("incorrect global version: %v, expected: %v", v, gv)
	}
}

func TestOutbox_ReadWrite(t *testing.T) {
	t.Parallel()

	o := RestoreOutbox(
		value.NewOutboxID(),
		value.NewPartitionID(2),
		value.NewGlobalVersion(12345))

	data := make([]byte, OutboxSize+2)
	data0 := o.WriteData(data)

	if l := len(data0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	o1 := &Outbox{}
	data1 := o1.ReadData(data)

	if l := len(data1); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	if *o1 != *o {
		t.Errorf("incorrect result: %v, expected: %v", *o1, *o)
	}
}
