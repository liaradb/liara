package link

import (
	"testing"

	"github.com/liaradb/liaradb/encoder/page"
)

func TestRecordID(t *testing.T) {
	fn := NewFileName("testfile")
	bid := NewBlockID(fn, 1)
	rid := NewRecordID(bid, 2)

	if b := rid.BlockID(); b != bid {
		t.Errorf("incorrect block id: %v, expected: %v", b, bid)
	}

	if p := rid.Position(); p != 2 {
		t.Errorf("incorrect position: %v, expected: %v", p, 2)
	}
}

func TestRecordID_Offset(t *testing.T) {
	fn := NewFileName("testfile")
	bid := NewBlockID(fn, 1)
	rid := NewRecordID(bid, 2)

	want := page.Offset(123 * 2)
	if o := rid.Offset(123); o != want {
		t.Errorf("incorrect offset: %v, expected: %v", o, want)
	}
}
