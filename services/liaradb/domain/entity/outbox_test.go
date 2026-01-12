package entity

import (
	"testing"

	"github.com/liaradb/liaradb/domain/value"
)

func TestOutbox_ReadWrite(t *testing.T) {
	o := RestoreOutbox(
		value.NewOutboxID(),
		value.NewPartitionRange(
			value.NewPartitionID(2),
			value.NewPartitionID(3)),
		value.NewGlobalVersion(12345))

	data := make([]byte, OutboxSize+2)
	data0 := o.Write(data)

	if l := len(data0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	o1 := &Outbox{}
	data1 := o1.Read(data)
	if l := len(data1); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	if *o1 != *o {
		t.Errorf("incorrect result: %v, expected: %v", *o1, *o)
	}
}
