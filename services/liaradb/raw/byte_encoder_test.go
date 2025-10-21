package raw

import (
	"slices"
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestByteEncoder(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	se := ByteEncoder{}

	want := []byte{1, 2, 3, 4, 5}
	if err := se.Write(w, want); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	result, err := se.Read(r)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(result, want) {
		t.Errorf("incorrect value: %v, expected: %v", result, want)
	}
}
