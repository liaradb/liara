package raw

import (
	"reflect"
	"testing"
)

func TestCopy_ExactLength(t *testing.T) {
	t.Parallel()

	a := []byte{1, 2, 3, 4, 5}
	b := make([]byte, 5)

	if err := Copy(b, a); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(a, b) {
		t.Error("copy failed")
	}
}

func TestCopy_DestinationTooShort(t *testing.T) {
	t.Parallel()

	a := []byte{1, 2, 3, 4, 5}
	b := make([]byte, 0, 5)

	if err := Copy(b, a); err == nil {
		t.Error("should be partial copy")
	}

	if reflect.DeepEqual(a, b) {
		t.Error("should not have fully copied")
	}
}

func TestCopy_DestinationNil(t *testing.T) {
	t.Parallel()

	a := []byte{1, 2, 3, 4, 5}
	var b []byte

	if err := Copy(b, a); err == nil {
		t.Error("should not copy")
	}
}

func TestCopy_SourceEmptyDestinationNil(t *testing.T) {
	t.Parallel()

	a := []byte{}
	var b []byte

	if err := Copy(b, a); err != nil {
		t.Error("should be noop")
	}
}

func TestCopy_SourceNil(t *testing.T) {
	t.Parallel()

	var a []byte
	b := make([]byte, 5)

	if err := Copy(b, a); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(b, make([]byte, 5)) {
		t.Error("should not have changed anything")
	}
}

func TestCopyAt_Underflow(t *testing.T) {
	t.Parallel()

	a := []byte{1, 2, 3, 4, 5}
	b := make([]byte, 6)

	if err := CopyAt(b, a, 1); err != nil {
		t.Error(err)
	}

	c := []byte{0, 1, 2, 3, 4, 5}
	if !reflect.DeepEqual(b, c) {
		t.Error("copy failed")
	}
}

func TestCopyAt_Overflow(t *testing.T) {
	t.Parallel()

	a := []byte{1, 2, 3, 4, 5}
	b := make([]byte, 5)

	if err := CopyAt(b, a, 1); err == nil {
		t.Error("should be partial copy")
	}
}

func TestCopyAt_BadOffset(t *testing.T) {
	t.Parallel()

	a := []byte{1, 2, 3, 4, 5}
	b := make([]byte, 5)

	if err := CopyAt(b, a, 6); err == nil {
		t.Error("should not copy")
	}

	if err := CopyAt(b, a, -1); err == nil {
		t.Error("should not copy")
	}
}
