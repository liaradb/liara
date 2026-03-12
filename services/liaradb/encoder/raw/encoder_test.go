package raw

import (
	"io"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/util/testing/testutil"
)

func TestStringSize(t *testing.T) {
	t.Parallel()

	value := "abcde"
	want := HeaderSize + len(value)

	s := StringSize(value)
	if s != want {
		t.Errorf("incorrect size: %v, expected: %v", s, want)
	}
}

func TestByteEncoder(t *testing.T) {
	t.Parallel()

	t.Run("should read from nil slice", func(t *testing.T) {
		r, w := testutil.NewReaderWriter()

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
	})

	t.Run("should read empty from nil slice", func(t *testing.T) {
		r, w := testutil.NewReaderWriter()

		want := []byte{}
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
	})

	t.Run("should read from small slice", func(t *testing.T) {
		r, w := testutil.NewReaderWriter()

		want := []byte{1, 2, 3, 4, 5}
		if err := Write(w, want); err != nil {
			t.Fatal(err)
		}

		if err := w.Flush(); err != nil {
			t.Fatal(err)
		}

		result := make([]byte, 1)
		if err := Read(r, &result); err != nil {
			t.Fatal(err)
		}

		if !slices.Equal(result, want) {
			t.Errorf("incorrect value: %v, expected: %v", result, want)
		}
	})

	t.Run("should read from sufficient slice", func(t *testing.T) {
		r, w := testutil.NewReaderWriter()

		want := []byte{1, 2, 3, 4, 5}
		if err := Write(w, want); err != nil {
			t.Fatal(err)
		}

		if err := w.Flush(); err != nil {
			t.Fatal(err)
		}

		result := make([]byte, 10)
		if err := Read(r, &result); err != nil {
			t.Fatal(err)
		}

		if !slices.Equal(result, want) {
			t.Errorf("incorrect value: %v, expected: %v", result, want)
		}
	})
}

func TestByteEncoder__Short(t *testing.T) {
	t.Parallel()

	t.Run("should return err on buffer too short for size", func(t *testing.T) {
		b := buffer.New(1)

		want := []byte{1, 2, 3, 4, 5}
		if err := Write(b, want); err != io.ErrShortWrite {
			t.Error("should return err")
		}

		if _, err := b.Seek(0, io.SeekStart); err != nil {
			t.Fatal(err)
		}

		var result []byte
		if err := Read(b, &result); err != io.EOF {
			t.Error("should return err")
		}
	})

	t.Run("should return err on buffer too short for value", func(t *testing.T) {
		b := buffer.New(5)

		want := []byte{1, 2, 3, 4, 5}
		if err := Write(b, want); err != io.ErrShortWrite {
			t.Error("should return err")
		}

		if _, err := b.Seek(0, io.SeekStart); err != nil {
			t.Fatal(err)
		}

		var result []byte
		if err := Read(b, &result); err != io.EOF {
			t.Error("should return err")
		}
	})
}

func TestStringEncoder(t *testing.T) {
	t.Parallel()

	r, w := testutil.NewReaderWriter()

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
		"should handle 1": {value: 1},
		"should handle 2": {value: 2},
		"should handle 3": {value: 3},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			r, w := testutil.NewReaderWriter()
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

func TestReadString__Error(t *testing.T) {
	t.Parallel()

	t.Run("should return error on empty reader", func(t *testing.T) {
		var s string
		b := buffer.New(0)

		if err := ReadString(b, &s); err == nil {
			t.Error("should return error")
		}
	})

	t.Run("should return error on short reader", func(t *testing.T) {
		var s string
		b := buffer.New(4)

		if err := WriteInt32(b, int32(1)); err != nil {
			t.Fatal(err)
		}

		if _, err := b.Seek(0, io.SeekStart); err != nil {
			t.Fatal(err)
		}

		if err := ReadString(b, &s); err == nil {
			t.Error("should return error")
		}
	})
}

func TestWriteString__Error(t *testing.T) {
	t.Parallel()

	t.Run("should return error on empty reader", func(t *testing.T) {
		var s string
		b := buffer.New(0)

		if err := WriteString(b, s); err == nil {
			t.Error("should return error")
		}
	})

	t.Run("should return error on short reader", func(t *testing.T) {
		var s string = "a"
		b := buffer.New(4)

		if err := WriteString(b, s); err == nil {
			t.Error("should return error")
		}
	})
}

func TestReadInt8__Error(t *testing.T) {
	var i int8
	b := buffer.New(0)
	if err := ReadInt8(b, &i); err == nil {
		t.Error("should return error")
	}
}

func TestReadInt16__Error(t *testing.T) {
	var i int16
	b := buffer.New(1)
	if err := ReadInt16(b, &i); err == nil {
		t.Error("should return error")
	}
}

func TestReadInt32__Error(t *testing.T) {
	var i int32
	b := buffer.New(1)
	if err := ReadInt32(b, &i); err == nil {
		t.Error("should return error")
	}
}

func TestReadInt64__Error(t *testing.T) {
	var i int64
	b := buffer.New(1)
	if err := ReadInt64(b, &i); err == nil {
		t.Error("should return error")
	}
}

func TestInt16(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value int16
	}{
		"should handle 0": {value: 0},
		"should handle 1": {value: 1},
		"should handle 2": {value: 2},
		"should handle 3": {value: 3},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			r, w := testutil.NewReaderWriter()
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
		"should handle 1": {value: 1},
		"should handle 2": {value: 2},
		"should handle 3": {value: 3},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			r, w := testutil.NewReaderWriter()
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
		"should handle 1": {value: 1},
		"should handle 2": {value: 2},
		"should handle 3": {value: 3},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			r, w := testutil.NewReaderWriter()
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
