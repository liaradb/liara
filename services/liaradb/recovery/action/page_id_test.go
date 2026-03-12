package action

import (
	"io"
	"testing"

	"github.com/liaradb/liaradb/util/testing/testutil"
)

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
		"should handle single size": {
			size:     10,
			pageSize: 10,
			id:       1},
		"should handle multiple size": {
			size:     20,
			pageSize: 10,
			id:       2},
		"should handle remainder size": {
			size:     22,
			pageSize: 10,
			id:       2},
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

func TestPageID_NewActivePageIDFromSize(t *testing.T) {
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
		"should handle single size": {
			size:     10,
			pageSize: 10,
			id:       0},
		"should handle multiple size": {
			size:     20,
			pageSize: 10,
			id:       1},
		"should handle remainder size": {
			size:     22,
			pageSize: 10,
			id:       2},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			id := NewActivePageIDFromSize(c.size, c.pageSize)
			if id != c.id {
				t.Errorf("%v: incorrect id: %v, expected: %v", message, id, c.id)
			}
		})
	}
}

func TestPageID_Size(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip     bool
		id       PageID
		pageSize int64
		size     int64
	}{
		"should handle zero pageSize": {
			id:       1,
			pageSize: 0,
			size:     0},
		"should handle zero id": {
			id:       0,
			pageSize: 10,
			size:     0},
		"should handle id 1": {
			id:       1,
			pageSize: 10,
			size:     10},
		"should handle id greater than 1": {
			id:       2,
			pageSize: 10,
			size:     20},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			s := c.id.Position(c.pageSize)
			if s != c.size {
				t.Errorf("%v: incorrect size: %v, expected: %v", message, s, c.size)
			}
		})
	}
}

func TestPageID_Write(t *testing.T) {
	t.Parallel()

	r, w := testutil.NewReaderWriter()

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
