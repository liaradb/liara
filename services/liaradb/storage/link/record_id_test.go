package link

import "testing"

func TestRecordID(t *testing.T) {
	fn := NewFileName("testfile")
	fbid := NewFileBlockID(fn, 1)
	rid := NewRecordID(fbid, 2)

	if b := rid.FileBlockID(); b != fbid {
		t.Errorf("incorrect block id: %v, expected: %v", b, fbid)
	}

	if s := rid.SlotID(); s != 2 {
		t.Errorf("incorrect slot id: %v, expected: %v", s, 2)
	}
}

func TestRecordID_Offset(t *testing.T) {
	fn := NewFileName("testfile")
	fbid := NewFileBlockID(fn, 1)
	rid := NewRecordID(fbid, 2)

	want := Offset(123 * 2)
	if o := rid.Offset(123); o != want {
		t.Errorf("incorrect offset: %v, expected: %v", o, want)
	}
}
