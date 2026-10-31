package link

import "testing"

func TestFileRecordID(t *testing.T) {
	fn := NewFileName("testfile")
	fbid := NewFileBlockID(fn, 1)
	frid := NewFileRecordID(fbid, 2)

	if b := frid.FileBlockID(); b != fbid {
		t.Errorf("incorrect block id: %v, expected: %v", b, fbid)
	}

	if s := frid.SlotID(); s != 2 {
		t.Errorf("incorrect slot id: %v, expected: %v", s, 2)
	}
}

func TestFileRecordID_Offset(t *testing.T) {
	fn := NewFileName("testfile")
	fbid := NewFileBlockID(fn, 1)
	frid := NewFileRecordID(fbid, 2)

	want := Offset(123 * 2)
	if o := frid.Offset(123); o != want {
		t.Errorf("incorrect offset: %v, expected: %v", o, want)
	}
}
