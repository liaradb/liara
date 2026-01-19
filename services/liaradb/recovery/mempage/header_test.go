package mempage

import (
	"testing"

	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/util/testutil"
)

func TestHeader(t *testing.T) {
	t.Parallel()

	r, w := testutil.NewReaderBuffer()
	pid := action.PageID(1)
	tlid := action.TimeLineID(2)
	rem := record.NewLength(3)

	h := newHeader(pid, tlid, rem)

	if err := h.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := h.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	h2 := header{}
	if err := h2.Read(r); err != nil {
		t.Fatal(err)
	}

	testHeader(t, h2, pid, tlid, rem)
}

func testHeader(
	t *testing.T,
	h header,
	pid action.PageID,
	tlid action.TimeLineID,
	rem record.Length,
) {
	t.Helper()
	testutil.Getter(t, h.ID, pid, "ID")
	testutil.Getter(t, h.TimeLineID, tlid, "TimeLineID")
	testutil.Getter(t, h.LengthRemaining, rem, "LengthRemaining")
}
