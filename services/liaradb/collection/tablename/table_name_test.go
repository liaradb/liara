package tablename

import "testing"

func TestTableName(t *testing.T) {
	t.Parallel()

	n := "testfile"
	tn := New(n)

	wantKV := "testfile.kv"
	if kv := tn.KeyValue(); kv != wantKV {
		t.Errorf("incorrect key value file: %v, expected: %v", kv, wantKV)
	}

	wantEL := "testfile.el"
	if el := tn.EventLog(); el != wantEL {
		t.Errorf("incorrect event log file: %v, expected: %v", el, wantEL)
	}

	wantIdx0 := "testfile--0.idx"
	if idx := tn.Index(0); idx != wantIdx0 {
		t.Errorf("incorrect index file: %v, expected: %v", idx, wantIdx0)
	}

	wantIdx1 := "testfile--1.idx"
	if idx := tn.Index(1); idx != wantIdx1 {
		t.Errorf("incorrect index file: %v, expected: %v", idx, wantIdx1)
	}
}
