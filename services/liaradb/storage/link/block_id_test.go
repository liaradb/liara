package link

import (
	"testing"
)

func TestBlockID(t *testing.T) {
	fn := NewFileName("testfile")
	pos := FilePosition(1)
	bid := NewBlockID(fn, pos)

	if n := bid.FileName(); n != fn {
		t.Errorf("incorrect file name: %v, expected: %v", n, fn)
	}

	if p := bid.Position(); p != pos {
		t.Errorf("incorrect position: %v, expected: %v", p, pos)
	}

	bid.SetPosition(2)
	if p := bid.Position(); p != 2 {
		t.Errorf("incorrect position: %v, expected: %v", p, 2)
	}
}

func TestBlockID_RecordID(t *testing.T) {
	fn := NewFileName("testfile")
	fPos := FilePosition(1)
	rSID := SlotID(2)
	bid := NewBlockID(fn, fPos)
	rid := NewRecordID(bid, rSID)

	if r := bid.RecordID(rSID); r != rid {
		t.Errorf("incorrect record id: %v, expected: %v", r, rid)
	}
}
