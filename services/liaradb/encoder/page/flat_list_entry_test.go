package page

import (
	"io"
	"testing"
)

func TestFlatListEntry(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	crc := NewCRC([]byte{1, 2, 3, 4, 5})

	var le = FlatListEntry{1, 12345, 67890, crc}
	if err := le.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := le.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var le2 FlatListEntry
	if err := le2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if le != le2 {
		t.Errorf("incorrect value: %v, expected: %v", le2, le)
	}
}
