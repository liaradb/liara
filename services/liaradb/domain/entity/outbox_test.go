package entity

import (
	"testing"

	"github.com/liaradb/liaradb/domain/value"
)

func TestOutbox(t *testing.T) {
	t.Parallel()

	oid := value.NewOutboxID()
	pr := value.NewPartitionRange(
		value.NewPartitionID(2),
		value.NewPartitionID(3))

	o := NewOutbox(oid, pr)

	if i := o.ID(); i != oid {
		t.Errorf("incorrect id: %v, expected: %v", i, oid)
	}

	if r := o.PartitionRange(); r != pr {
		t.Errorf("incorrect partition range: %v, expected: %v", r, pr)
	}

	if v := o.GlobalVersion().Value(); v != 0 {
		t.Errorf("incorrect global version: %v, expected: %v", v, 0)
	}
}

func TestOutbox_UpdateGlobalVersion(t *testing.T) {
	t.Parallel()

	oid := value.NewOutboxID()
	pr := value.NewPartitionRange(
		value.NewPartitionID(2),
		value.NewPartitionID(3))

	o := NewOutbox(oid, pr)

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
		value.NewPartitionRange(
			value.NewPartitionID(2),
			value.NewPartitionID(3)),
		value.NewGlobalVersion(12345))

	data := make([]byte, OutboxSize+2)
	data0, ok := o.Write(data)
	if !ok {
		t.Error("unable to write")
	}

	if l := len(data0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	o1 := &Outbox{}
	data1, ok := o1.Read(data)
	if !ok {
		t.Error("unable to read")
	}
	if l := len(data1); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	if *o1 != *o {
		t.Errorf("incorrect result: %v, expected: %v", *o1, *o)
	}
}
