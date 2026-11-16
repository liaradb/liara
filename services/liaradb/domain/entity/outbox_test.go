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

	if s := o.Size(); s != OutboxSize {
		t.Errorf("incorrect size: %v, expected: %v", s, OutboxSize)
	}

	if i := o.ID(); i != oid {
		t.Errorf("incorrect id: %v, expected: %v", i, oid)
	}

	if p := o.PartitionID(); p != pid {
		t.Errorf("incorrect partition id: %v, expected: %v", p, pid)
	}

	if v := o.GlobalVersion().Value(); v != 0 {
		t.Errorf("incorrect global version: %v, expected: %v", v, 0)
	}

	if v := o.SubVersion().Value(); v != 0 {
		t.Errorf("incorrect sub version: %v, expected: %v", v, 0)
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
	sv := value.NewSubVersion(67890)

	o.UpdatePosition(gv, sv)

	if v := o.GlobalVersion(); v != gv {
		t.Errorf("incorrect global version: %v, expected: %v", v, gv)
	}

	if v := o.SubVersion(); v != sv {
		t.Errorf("incorrect global version: %v, expected: %v", v, sv)
	}
}

func TestOutbox_ReadWrite(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	o := RestoreOutbox(
		value.NewOutboxID(),
		value.NewPartitionID(2),
		value.NewGlobalVersion(12345),
		value.NewSubVersion(67890))

	if err := o.Write(w); err != nil {
		t.Fatal(err)
	}

	o1 := &Outbox{}
	if err := o1.Read(r); err != nil {
		t.Fatal(err)
	}

	if *o1 != *o {
		t.Errorf("incorrect result: %v, expected: %v", *o1, *o)
	}
}
