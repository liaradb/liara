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

	var lsn = Row{
		Data: value.NewData([]byte{}),
	}
	if err := lsn.Write(w); err != nil {
		t.Fatal(err)
	}

	// size := w.Len()
	// if s := lsn.Size(); s != size {
	// 	t.Errorf("incorrect size: %v, expected: %v", s, size)
	// }

	var lsn2 Row
	if err := lsn2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	// TODO: Create another comparison
	// Data comparison doesn't allow nil slice
	if !reflect.DeepEqual(lsn, lsn2) {
		t.Errorf("incorrect value: %v, expected: %v", lsn2, lsn)
	}
}
