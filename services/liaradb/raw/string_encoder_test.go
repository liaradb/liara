package raw

import (
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestStringEncoder(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	se := StringEncoder{}

	want := "abcde"
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

	if result != want {
		t.Errorf("incorrect value: %v, expected: %v", result, want)
	}
}
