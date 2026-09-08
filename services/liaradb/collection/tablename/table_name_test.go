package tablename

import (
	"testing"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage/link"
)

func TestTableName_Default(t *testing.T) {
	t.Parallel()

	tn := NewFromString("")

	wantKV := link.NewFileName("00000000-0000-0000-0000-000000000000--00000001.kv")
	if kv := tn.KeyValue(value.NewPartitionID(1)); kv != wantKV {
		t.Errorf("incorrect key value file: %v, expected: %v", kv, wantKV)
	}

	wantEL := link.NewFileName("00000000-0000-0000-0000-000000000000--00000001.el")
	if el := tn.EventLog(value.NewPartitionID(1)); el != wantEL {
		t.Errorf("incorrect event log file: %v, expected: %v", el, wantEL)
	}

	wantIdx0 := link.NewFileName("00000000-0000-0000-0000-000000000000--00000000--00000002.idx")
	if idx := tn.Index(0, value.NewPartitionID(2)); idx != wantIdx0 {
		t.Errorf("incorrect index file: %v, expected: %v", idx, wantIdx0)
	}

	wantIdx1 := link.NewFileName("00000000-0000-0000-0000-000000000000--00000001--00000002.idx")
	if idx := tn.Index(1, value.NewPartitionID(2)); idx != wantIdx1 {
		t.Errorf("incorrect index file: %v, expected: %v", idx, wantIdx1)
	}
}

func TestTableName(t *testing.T) {
	t.Parallel()

	tn := NewFromString("testfile")

	wantKV := link.NewFileName("00000000-0000-0000-0000-000000000000--testfile--00000001.kv")
	if kv := tn.KeyValue(value.NewPartitionID(1)); kv != wantKV {
		t.Errorf("incorrect key value file: %v, expected: %v", kv, wantKV)
	}

	wantEL := link.NewFileName("00000000-0000-0000-0000-000000000000--testfile--00000001.el")
	if el := tn.EventLog(value.NewPartitionID(1)); el != wantEL {
		t.Errorf("incorrect event log file: %v, expected: %v", el, wantEL)
	}

	wantIdx0 := link.NewFileName("00000000-0000-0000-0000-000000000000--testfile--00000000--00000002.idx")
	if idx := tn.Index(0, value.NewPartitionID(2)); idx != wantIdx0 {
		t.Errorf("incorrect index file: %v, expected: %v", idx, wantIdx0)
	}

	wantIdx1 := link.NewFileName("00000000-0000-0000-0000-000000000000--testfile--00000001--00000002.idx")
	if idx := tn.Index(1, value.NewPartitionID(2)); idx != wantIdx1 {
		t.Errorf("incorrect index file: %v, expected: %v", idx, wantIdx1)
	}
}
