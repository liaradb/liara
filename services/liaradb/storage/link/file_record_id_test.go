package link

import "testing"

func TestFileRecordID(t *testing.T) {
	fn := NewFileName("testfile")
	fbid := NewFileBlockID(fn, 1)
	frid := NewFileRecordID(fbid, 2)

	if s, want := frid.Size(), fn.Size()+8+2; s != want {
		t.Errorf("incorrect size: %v, expected: %v", s, want)
	}

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

func TestFileRecordID_ReadWrite(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	frid0 := NewFileRecordID(NewFileBlockID(NewFileName("testfile"), 2), 3)
	if err := frid0.Write(w); err != nil {
		t.Fatal(err)
	}

	var frid1 FileRecordID
	if err := frid1.Read(r); err != nil {
		t.Fatal(err)
	}

	if frid1 != frid0 {
		t.Errorf("incorrect result: %v, expected: %v", frid1, frid0)
	}
}
