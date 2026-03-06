package raw

import (
	"io"
	"slices"
	"testing"
)

var (
	data0 = []byte{1, 2, 3, 4, 5}
	data1 = []byte{6, 7, 8, 9, 0}
)

func TestBuffer_Default(t *testing.T) {
	t.Parallel()

	b := NewBuffer(10)

	result := make([]byte, 10)
	if n, err := b.Read(result); err != nil {
		t.Error(err)
	} else if n != 10 {
		t.Errorf("incorrect count: %v, expected: %v", n, 10)
	}

	empty := make([]byte, 10)
	if !slices.Equal(result, empty) {
		t.Errorf("incorrect result: %v, expected: %v", result, empty)
	}
}

func TestBuffer_NewBufferFromSlice(t *testing.T) {
	t.Parallel()

	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	b := NewBufferFromSlice(data)

	result := make([]byte, 10)
	if n, err := b.Read(result); err != nil {
		t.Error(err)
	} else if n != 10 {
		t.Errorf("incorrect count: %v, expected: %v", n, 10)
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

func TestBuffer_Bytes(t *testing.T) {
	t.Parallel()

	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	b := NewBufferFromSlice(data)

	if b := b.Bytes(); !slices.Equal(b, data) {
		t.Errorf("incorrect byte slice: %v, expected: %v", b, data)
	}
}

func TestBuffer_Length(t *testing.T) {
	t.Parallel()

	var size int64 = 10
	b := NewBuffer(size)

	if s := b.Length(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", size, s)
	}
}

func TestBuffer_Clear(t *testing.T) {
	t.Parallel()

	b := NewBuffer(10)

	if n, err := b.Write(data0); err != nil {
		t.Error(err)
	} else if n != 5 {
		t.Errorf("incorrect count: %v, expected: %v", n, 5)
	}

	// Read next section
	result := make([]byte, 5)
	if n, err := b.Read(result); err != nil {
		t.Error(err)
	} else if n != 5 {
		t.Errorf("incorrect count: %v, expected: %v", n, 10)
	}

	empty := make([]byte, 5)
	if !slices.Equal(result, empty) {
		t.Errorf("incorrect result: %v, expected: %v", result, empty)
	}

	b.Clear()

	// Read first section
	if n, err := b.ReadAt(result, 0); err != nil {
		t.Error(err)
	} else if n != 5 {
		t.Errorf("incorrect count: %v, expected: %v", n, 10)
	}

	if !slices.Equal(result, empty) {
		t.Errorf("incorrect result: %v, expected: %v", result, empty)
	}
}

func TestBuffer_ClearAfter(t *testing.T) {
	t.Parallel()

	base := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	b := NewBufferFromSlice(slices.Clone(base))

	if !slices.Equal(b.Bytes(), base) {
		t.Errorf("incorrect bytes: %v, expected: %v", b.Bytes(), base)
	}

	b.ClearAfter(4)

	want := []byte{1, 2, 3, 4, 0, 0, 0, 0}
	if !slices.Equal(b.Bytes(), want) {
		t.Errorf("incorrect bytes: %v, expected: %v", b.Bytes(), want)
	}
}

func TestBuffer_Reset(t *testing.T) {
	t.Parallel()

	b0 := NewBuffer(20)

	if n, err := b0.Write(data0); err != nil {
		t.Error(err)
	} else if n != 5 {
		t.Errorf("incorrect count: %v, expected: %v", n, 5)
	}

	b1 := NewBuffer(20)

	if n, err := b1.Write(data1); err != nil {
		t.Error(err)
	} else if n != 5 {
		t.Errorf("incorrect count: %v, expected: %v", n, 5)
	}

	b1.Reset(b0.Bytes())

	result := make([]byte, 5)

	if n, err := b1.Read(result); err != nil {
		t.Error(err)
	} else if n != 5 {
		t.Errorf("incorrect count: %v, expected: %v", n, 5)
	} else if !slices.Equal(result, data0) {
		t.Errorf("incorrect result: %v, expected: %v", result, data0)
	}
}

func TestBuffer_ReadWrite(t *testing.T) {
	t.Parallel()

	t.Run("should write", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(20)

		if n, err := b.Write(data0); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		}

		if n, err := b.Write(data1); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		}

		if n, err := b.Seek(0, io.SeekStart); err != nil {
			t.Error(err)
		} else if n != 0 {
			t.Errorf("incorrect count: %v, expected: %v", n, 0)
		}

		result := make([]byte, 5)

		if n, err := b.Read(result); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		} else if !slices.Equal(result, data0) {
			t.Errorf("incorrect result: %v, expected: %v", result, data0)
		}

		if n, err := b.Read(result); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		} else if !slices.Equal(result, data1) {
			t.Errorf("incorrect result: %v, expected: %v", result, data0)
		}
	})

	t.Run("should not write after buffer", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(2)

		result := make([]byte, 5)

		n, err := b.Write(result)
		if err != io.ErrShortWrite {
			t.Error("should return ErrShortWrite")
		}
		if n != 2 {
			t.Errorf("should return remainder: %v, expected: %v", n, 2)
		}
	})

	t.Run("should not read after buffer", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(2)

		result := make([]byte, 5)

		n, err := b.Read(result)
		if err != io.EOF {
			t.Error("should return EOF")
		}
		if n != 2 {
			t.Errorf("should return remainder: %v, expected: %v", n, 2)
		}
	})
}

func TestBuffer_ReadAtWriteAt(t *testing.T) {
	t.Parallel()

	t.Run("should write at", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(256)

		// Write out of order
		if n, err := b.WriteAt(data1, 5); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		}

		if n, err := b.WriteAt(data0, 0); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		}

		result := make([]byte, 5)

		// Read out of order
		if n, err := b.ReadAt(result, 5); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		} else if !slices.Equal(result, data1) {
			t.Errorf("incorrect result: %v, expected: %v", result, data1)
		}

		if n, err := b.ReadAt(result, 0); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		} else if !slices.Equal(result, data0) {
			t.Errorf("incorrect result: %v, expected: %v", result, data0)
		}
	})

	t.Run("should update cursor position", func(t *testing.T) {
		b := NewBuffer(256)

		if n, err := b.WriteAt(data0, 0); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		}

		if n, err := b.Write(data1); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		}

		if n, err := b.Seek(0, io.SeekStart); err != nil {
			t.Error(err)
		} else if n != 0 {
			t.Errorf("incorrect count: %v, expected: %v", n, 0)
		}

		result := make([]byte, 5)

		if n, err := b.Read(result); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		} else if !slices.Equal(result, data0) {
			t.Errorf("incorrect result: %v, expected: %v", result, data0)
		}

		if n, err := b.Read(result); err != nil {
			t.Error(err)
		} else if n != 5 {
			t.Errorf("incorrect count: %v, expected: %v", n, 5)
		} else if !slices.Equal(result, data1) {
			t.Errorf("incorrect result: %v, expected: %v", result, data1)
		}
	})

	t.Run("should not write before buffer", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(20)

		result := make([]byte, 5)

		n, err := b.WriteAt(result, -2)
		if err != ErrUnderflow {
			t.Error("should return underflow error")
		}
		if n != 0 {
			t.Error("should return 0")
		}
	})

	t.Run("should not write beyond buffer", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(20)

		result := make([]byte, 5)

		n, err := b.WriteAt(result, 18)
		if err != io.ErrShortWrite {
			t.Error("should return ErrShortWrite")
		}
		if n != 2 {
			t.Errorf("should return remainder: %v, expected: %v", n, 2)
		}
	})

	t.Run("should not write after buffer", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(20)

		result := make([]byte, 5)

		n, err := b.WriteAt(result, 22)
		if err != io.ErrShortWrite {
			t.Error("should return ErrShortWrite")
		}
		if n != 0 {
			t.Errorf("should return remainder: %v, expected: %v", n, 0)
		}
	})

	t.Run("should not read before buffer", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(20)

		result := make([]byte, 5)

		n, err := b.ReadAt(result, -2)
		if err != ErrUnderflow {
			t.Error("should return underflow error")
		}
		if n != 0 {
			t.Error("should return 0")
		}
	})

	t.Run("should not read beyond buffer", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(20)

		result := make([]byte, 5)

		n, err := b.ReadAt(result, 18)
		if err != io.EOF {
			t.Error("should return EOF")
		}
		if n != 2 {
			t.Errorf("should return remainder: %v, expected: %v", n, 2)
		}
	})

	t.Run("should not read after buffer", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(20)

		result := make([]byte, 5)

		n, err := b.ReadAt(result, 22)
		if err != io.EOF {
			t.Error("should return EOF")
		}
		if n != 0 {
			t.Errorf("should return remainder: %v, expected: %v", n, 0)
		}
	})

	t.Run("should not read at buffer end", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(20)

		result := make([]byte, 5)

		n, err := b.ReadAt(result, 20)
		if err != io.EOF {
			t.Error("should return EOF")
		}
		if n != 0 {
			t.Errorf("should return remainder: %v, expected: %v", n, 0)
		}
	})

	t.Run("should read trivial case", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(20)

		result := make([]byte, 0)

		n, err := b.ReadAt(result, 0)
		if err != nil {
			t.Error(err)
		}
		if n != 0 {
			t.Errorf("should return remainder: %v, expected: %v", n, 0)
		}
	})

	t.Run("should read trivial case at buffer end", func(t *testing.T) {
		t.Parallel()

		b := NewBuffer(20)

		result := make([]byte, 0)

		n, err := b.ReadAt(result, 20)
		if err != nil {
			t.Error(err)
		}
		if n != 0 {
			t.Errorf("should return remainder: %v, expected: %v", n, 0)
		}
	})
}

func TestBuffer_Seek(t *testing.T) {
	t.Parallel()

	const initialPosition = 10

	for message, c := range map[string]struct {
		skip     bool
		position int64
		whence   int
		err      error
		n        int
	}{
		"should handle defaults": {
			position: 0,
			whence:   0,
			err:      nil,
			n:        0},
		"should not seek to negative position from start": {
			position: -1,
			whence:   io.SeekStart,
			err:      ErrUnderflow,
			n:        initialPosition},
		"should not seek to negative position from current": {
			position: -11,
			whence:   io.SeekCurrent,
			err:      ErrUnderflow,
			n:        initialPosition},
		"should not seek to negative position from end": {
			position: -21,
			whence:   io.SeekEnd,
			err:      ErrUnderflow,
			n:        initialPosition},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			b := NewBuffer(20)

			if n, err := b.Seek(initialPosition, io.SeekStart); err != nil {
				t.Error(err)
			} else if n != initialPosition {
				t.Errorf("%v: incorrect count: %v, expected: %v", message, n, initialPosition)
			}

			n, err := b.Seek(c.position, c.whence)
			if err != c.err {
				if c.err == nil {
					t.Error(err)
				} else {
					t.Errorf("%v: incorrect error: %v, expected: %v", message, err, c.err)
				}
			}
			if n != int64(c.n) {
				t.Errorf("%v: incorrect n: %v, expected: %v", message, n, c.n)
			}
		})
	}
}
