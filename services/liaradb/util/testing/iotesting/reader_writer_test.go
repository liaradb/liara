package iotesting

import (
	"slices"
	"testing"
)

func TestNewReaderWriter(t *testing.T) {
	t.Parallel()

	r, w := NewReaderWriter()

	want := []byte("abcde")
	if n, err := w.Write(want); err != nil {
		t.Fatal(err)
	} else if n != 5 {
		t.Errorf("incorrect length: %v, expected: %v", n, 5)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	result := make([]byte, 5)
	if n, err := r.Read(result); err != nil {
		t.Fatal(err)
	} else if n != 5 {
		t.Errorf("incorrect length: %v, expected: %v", n, 5)
	}

	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func TestNewReaderBuffer(t *testing.T) {
	t.Parallel()

	r, w := NewReaderBuffer()

	want := []byte("abcde")
	if n, err := w.Write(want); err != nil {
		t.Fatal(err)
	} else if n != 5 {
		t.Errorf("incorrect length: %v, expected: %v", n, 5)
	}

	result := make([]byte, 5)
	if n, err := r.Read(result); err != nil {
		t.Fatal(err)
	} else if n != 5 {
		t.Errorf("incorrect length: %v, expected: %v", n, 5)
	}

	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}
