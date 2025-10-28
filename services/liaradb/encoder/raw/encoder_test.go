package raw

import (
	"slices"
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestStringSize(t *testing.T) {
	value := "abcde"
	want := HeaderSize + len(value)

	s := StringSize(value)
	if s != want {
		t.Errorf("incorrect size: %v, expected: %v", s, want)
	}
}

func TestByteEncoder(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	want := []byte{1, 2, 3, 4, 5}
	if err := Write(w, want); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var result []byte
	if err := Read(r, &result); err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(result, want) {
		t.Errorf("incorrect value: %v, expected: %v", result, want)
	}
}

func TestStringEncoder(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	want := "abcde"
	if err := WriteString(w, want); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var result string
	if err := ReadString(r, &result); err != nil {
		t.Fatal(err)
	}

	if result != want {
		t.Errorf("incorrect value: %v, expected: %v", result, want)
	}
}
