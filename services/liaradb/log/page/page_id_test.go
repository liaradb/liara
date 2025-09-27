package page

import (
	"io"
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestPageID(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	var pid PageID = 123456
	if err := pid.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var pid2 PageID
	if err := pid2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if pid != pid2 {
		t.Errorf("incorrect value: %v, expected: %v", pid2, pid)
	}
}

func TestPageID_NewPageIDFromSize(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip     bool
		size     int64
		pageSize int64
		id       PageID
	}{
		"should handle zero pageSize": {
			size:     10,
			pageSize: 0,
			id:       0},
		"should handle zero size": {
			size:     0,
			pageSize: 10,
			id:       0},
		"should handle multiple size": {
			size:     20,
			pageSize: 10,
			id:       2},
		"should handle remainder size": {
			size:     22,
			pageSize: 10,
			id:       3},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			id := NewPageIDFromSize(c.size, c.pageSize)
			if id != c.id {
				t.Errorf("%v: incorrect id: %v, expected: %v", message, id, c.id)
			}
		})
	}
}
