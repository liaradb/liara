package page

import (
	"io"
	"testing"

	"github.com/liaradb/liaradb/util/testutil"
)

func TestMagic(t *testing.T) {
	t.Parallel()

	r, w := testutil.NewReaderWriter()

	var m Magic = MagicPage
	if err := m.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var m2 Magic
	if err := m2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if m != m2 {
		t.Errorf("incorrect value: %v, expected: %v", m2, m)
	}
}

func TestMagicPage(t *testing.T) {
	t.Parallel()

	if s := MagicPage.String(); s != "PAGE" {
		t.Error("value is incorrect")
	}
}
