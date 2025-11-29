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

func TestInt8(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value int8
	}{
		"should handle 0": {value: 0},
		"should handle 1": {value: 0},
		"should handle 2": {value: 0},
		"should handle 3": {value: 0},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			r, w := assert.NewReaderWriter()
			var want, result int8 = c.value, 0

			if err := WriteInt8(w, want); err != nil {
				t.Fatal(err)
			}

			if err := w.Flush(); err != nil {
				t.Fatal(err)
			}

			if err := ReadInt8(r, &result); err != nil {
				t.Fatal(err)
			}

			if result != want {
				t.Errorf("incorrect result: %v, expected: %v", result, want)
			}
		})
	}
}

func TestInt16(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value int16
	}{
		"should handle 0": {value: 0},
		"should handle 1": {value: 0},
		"should handle 2": {value: 0},
		"should handle 3": {value: 0},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			r, w := assert.NewReaderWriter()
			var want, result int16 = c.value, 0

			if err := WriteInt16(w, want); err != nil {
				t.Fatal(err)
			}

			if err := w.Flush(); err != nil {
				t.Fatal(err)
			}

			if err := ReadInt16(r, &result); err != nil {
				t.Fatal(err)
			}

			if result != want {
				t.Errorf("incorrect result: %v, expected: %v", result, want)
			}
		})
	}
}

func TestInt32(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value int32
	}{
		"should handle 0": {value: 0},
		"should handle 1": {value: 0},
		"should handle 2": {value: 0},
		"should handle 3": {value: 0},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			r, w := assert.NewReaderWriter()
			var want, result int32 = c.value, 0

			if err := WriteInt32(w, want); err != nil {
				t.Fatal(err)
			}

			if err := w.Flush(); err != nil {
				t.Fatal(err)
			}

			if err := ReadInt32(r, &result); err != nil {
				t.Fatal(err)
			}

			if result != want {
				t.Errorf("incorrect result: %v, expected: %v", result, want)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value int64
	}{
		"should handle 0": {value: 0},
		"should handle 1": {value: 0},
		"should handle 2": {value: 0},
		"should handle 3": {value: 0},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			r, w := assert.NewReaderWriter()
			var want, result int64 = c.value, 0

			if err := WriteInt64(w, want); err != nil {
				t.Fatal(err)
			}

			if err := w.Flush(); err != nil {
				t.Fatal(err)
			}

			if err := ReadInt64(r, &result); err != nil {
				t.Fatal(err)
			}

			if result != want {
				t.Errorf("incorrect result: %v, expected: %v", result, want)
			}
		})
	}
}
