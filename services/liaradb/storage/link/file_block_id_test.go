package link

import "testing"

func TestFileBlockID(t *testing.T) {
	t.Parallel()

	fn := NewFileName("testfile")
	pos := BlockID(1)
	fbid := NewFileBlockID(fn, pos)

	if s, want := fbid.Size(), fn.Size()+8; s != want {
		t.Errorf("incorrect size: %v, expected: %v", s, want)
	}

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
	t.Parallel()

	fn := NewFileName("testfile")
	fPos := BlockID(1)
	rSID := SlotID(2)
	fbid := NewFileBlockID(fn, fPos)
	frid := NewFileRecordID(fbid, rSID)

	if r := fbid.FileRecordID(rSID); r != frid {
		t.Errorf("incorrect record id: %v, expected: %v", r, frid)
	}
}

func TestFileBlockID_Next(t *testing.T) {
	t.Parallel()

	fbid0 := NewFileBlockID(NewFileName("testfile"), 2)
	fbid1 := fbid0.Next()

	if fn1, fn0 := fbid1.FileName(), fbid0.FileName(); fn1 != fn0 {
		t.Errorf("incorrect file name: %v, expected: %v", fn1, fn0)
	}

	if bid, want := fbid1.BlockID(), BlockID(3); bid != want {
		t.Errorf("incorrect block id: %v, expected: %v", bid, want)
	}
}

func TestFileBlockID_RecordID(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip bool
		fbid FileBlockID
		sid  SlotID
		rid  RecordID
	}{
		"should handle 0": {
			fbid: NewFileBlockID(NewFileName("testfile"), 2),
			sid:  0,
			rid:  NewRecordID(2, 0),
		},
		"should handle 1": {
			fbid: NewFileBlockID(NewFileName("testfile"), 2),
			sid:  1,
			rid:  NewRecordID(2, 1),
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			rid := c.fbid.RecordID(c.sid)
			want := NewRecordID(c.fbid.BlockID(), c.sid)
			if rid != want {
				t.Errorf("incorrect record id: %v, expected: %v", rid, want)
			}
		})
	}
}

func TestFileBlockID_ReadWrite(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	fbid0 := NewFileBlockID(NewFileName("testfile"), 2)
	if err := fbid0.Write(w); err != nil {
		t.Fatal(err)
	}

	var fbid1 FileBlockID
	if err := fbid1.Read(r); err != nil {
		t.Fatal(err)
	}

	if fbid1 != fbid0 {
		t.Errorf("incorrect result: %v, expected: %v", fbid1, fbid0)
	}
}
