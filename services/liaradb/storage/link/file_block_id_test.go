package link

import "testing"

func TestFileBlockID(t *testing.T) {
	fn := NewFileName("testfile")
	pos := BlockID(1)
	fbid := NewFileBlockID(fn, pos)

	if n := fbid.FileName(); n != fn {
		t.Errorf("incorrect file name: %v, expected: %v", n, fn)
	}

	if p := fbid.BlockID(); p != pos {
		t.Errorf("incorrect block id: %v, expected: %v", p, pos)
	}

	fbid.SetBlockID(2)
	if p := fbid.BlockID(); p != 2 {
		t.Errorf("incorrect block id: %v, expected: %v", p, 2)
	}
}

func TestFileBlockID_FileRecordID(t *testing.T) {
	fn := NewFileName("testfile")
	fPos := BlockID(1)
	rSID := SlotID(2)
	fbid := NewFileBlockID(fn, fPos)
	frid := NewFileRecordID(fbid, rSID)

	if r := fbid.FileRecordID(rSID); r != frid {
		t.Errorf("incorrect record id: %v, expected: %v", r, frid)
	}
}
