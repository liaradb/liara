package action

import (
	"io"
	"testing"

	"github.com/liaradb/liaradb/util/testutil"
)

func TestTimeLineID(t *testing.T) {
	t.Parallel()

	r, w := testutil.NewReaderWriter()

	var tlid TimeLineID = 123456
	if err := tlid.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var tlid2 TimeLineID
	if err := tlid2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if tlid != tlid2 {
		t.Errorf("incorrect value: %v, expected: %v", tlid2, tlid)
	}
}
