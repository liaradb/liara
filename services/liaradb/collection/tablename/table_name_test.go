package tablename

import (
	"testing"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/storage/link"
)

func TestTableName(t *testing.T) {
	t.Parallel()

	n := value.TenantID("testfile")
	tn := New(n)

	wantKV := link.NewFileName("testfile--1.kv")
	if kv := tn.KeyValue(value.NewPartitionID(1)); kv != wantKV {
		t.Errorf("incorrect key value file: %v, expected: %v", kv, wantKV)
	}

	wantEL := link.NewFileName("testfile--1.el")
	if el := tn.EventLog(value.NewPartitionID(1)); el != wantEL {
		t.Errorf("incorrect event log file: %v, expected: %v", el, wantEL)
	}

	wantIdx0 := link.NewFileName("testfile--0--2.idx")
	if idx := tn.Index(0, value.NewPartitionID(2)); idx != wantIdx0 {
		t.Errorf("incorrect index file: %v, expected: %v", idx, wantIdx0)
	}

	wantIdx1 := link.NewFileName("testfile--1--2.idx")
	if idx := tn.Index(1, value.NewPartitionID(2)); idx != wantIdx1 {
		t.Errorf("incorrect index file: %v, expected: %v", idx, wantIdx1)
	}
}
