package entity

import (
	"io"
	"reflect"
	"testing"

	"github.com/liaradb/liaradb/domain/value"
)

func TestRow(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var row = Row{}
	row.SetData(value.NewData([]byte{}))

	if err := row.Write(w); err != nil {
		t.Fatal(err)
	}

	// size := w.Len()
	// if s := lsn.Size(); s != size {
	// 	t.Errorf("incorrect size: %v, expected: %v", s, size)
	// }

	var row2 Row
	if err := row2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	// TODO: Create another comparison
	// Data comparison doesn't allow nil slice
	if !reflect.DeepEqual(row, row2) {
		t.Errorf("incorrect value: %v, expected: %v", row2, row)
	}
}
