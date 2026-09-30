package int16list

import (
	"slices"
	"testing"

	"github.com/liaradb/liaradb/util/testing/should"
)

func TestInt16List_Default(t *testing.T) {
	t.Parallel()

	l := New([]byte{})

	if length := l.Length(); length != 0 {
		t.Errorf("incorrect length: %v, expected: %v", length, 0)
	}

	if s := l.Size(); s != 0 {
		t.Errorf("incorrect size: %v, expected: %v", s, 0)
	}
}

func TestList_GetSet(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 16))

	if length := l.Length(); length != 16 {
		t.Errorf("incorrect length: %v, expected: %v", length, 16)
	}

	if s := l.Size(); s != 8 {
		t.Errorf("incorrect size: %v, expected: %v", s, 1)
	}

	for i := range int16(8) {
		l.Set(i, (i+1)*11)
	}

	should.Panic(t, "should not set value beyond size", func() {
		l.Set(8, 55)
	})

	for i := range int16(8) {
		want := (i + 1) * 11
		if v := l.Get(i); v != want {
			t.Errorf("incorrect value: %v, expected: %v", v, want)
		}
	}

	should.Panic(t, "should not set value beyond size", func() {
		_ = l.Get(8)
	})
}

func TestList_GetInt32SetInt32(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 16))

	if length := l.Length(); length != 16 {
		t.Errorf("incorrect length: %v, expected: %v", length, 16)
	}

	if s := l.Size(); s != 8 {
		t.Errorf("incorrect size: %v, expected: %v", s, 1)
	}

	for i := range int16(4) {
		l.SetInt32(i*2, int32((i+1)*11))
	}

	should.Panic(t, "should not set value beyond size", func() {
		l.SetInt32(8, 55)
	})

	for i := range int16(4) {
		want := int32((i + 1) * 11)
		if v := l.GetInt32(i * 2); v != want {
			t.Errorf("incorrect value: %v, expected: %v", v, want)
		}
	}

	should.Panic(t, "should not set value beyond size", func() {
		_ = l.GetInt32(8)
	})
}

func TestList_GetInt64SetInt64(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 16))

	if length := l.Length(); length != 16 {
		t.Errorf("incorrect length: %v, expected: %v", length, 16)
	}

	if s := l.Size(); s != 8 {
		t.Errorf("incorrect size: %v, expected: %v", s, 1)
	}

	for i := range int16(2) {
		l.SetInt64(i*4, int64((i+1)*11))
	}

	should.Panic(t, "should not set value beyond size", func() {
		l.SetInt64(8, 55)
	})

	for i := range int16(2) {
		want := int64((i + 1) * 11)
		if v := l.GetInt64(i * 4); v != want {
			t.Errorf("incorrect value: %v, expected: %v", v, want)
		}
	}

	should.Panic(t, "should not set value beyond size", func() {
		_ = l.GetInt64(8)
	})
}

func TestInt16List_Shift(t *testing.T) {
	t.Parallel()

	data := []int16{1, 2, 3, 4, 5, 6, 7, 8}

	for message, c := range map[string]struct {
		skip    bool
		want    []int16
		index   int16
		count   int16
		succeed bool
	}{
		"should not shift negative index": {
			want:    data,
			index:   -2,
			count:   2,
			succeed: false,
		},
		"should not shift negative count": {
			want:    data,
			index:   2,
			count:   -2,
			succeed: false,
		},
		"should shift count 0": {
			want:    []int16{1, 2, 3, 4, 5, 6, 7, 8},
			index:   2,
			count:   0,
			succeed: true,
		},
		"should shift": {
			want:    []int16{1, 2, 1, 2, 3, 4, 5, 6},
			index:   2,
			count:   2,
			succeed: true,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			l := New(make([]byte, 32))

			for i, d := range data {
				l.Set(int16(i), d)
			}

			if c.succeed {
				l.Shift(c.index, c.count)
			} else {
				should.Panic(t, "should not shift", func() {
					l.Shift(c.index, c.count)
				})
			}

			result := make([]int16, 0, len(data))
			for i := range data {
				d := l.Get(int16(i))
				result = append(result, d)
			}

			if !slices.Equal(result, c.want) {
				t.Errorf("incorrect : %v, expected: %v", result, c.want)
			}
		})
	}
}

func TestInt16List_ShiftRange(t *testing.T) {
	t.Parallel()

	data := []int16{1, 2, 3, 4, 5, 6, 7, 8}

	for message, c := range map[string]struct {
		skip    bool
		want    []int16
		start   int16
		end     int16
		shift   int16
		succeed bool
	}{
		"should not shift negative start": {
			want:    data,
			start:   -2,
			end:     2,
			shift:   2,
			succeed: false,
		},
		"should not shift negative length": {
			want:    data,
			start:   2,
			end:     0,
			shift:   2,
			succeed: false,
		},
		"should shift count 0": {
			want:    []int16{1, 2, 3, 4, 5, 6, 7, 8},
			start:   2,
			end:     4,
			shift:   0,
			succeed: true,
		},
		"should shift length 0": {
			want:    []int16{1, 2, 3, 4, 5, 6, 7, 8},
			start:   2,
			end:     2,
			shift:   2,
			succeed: true,
		},
		"should shift left": {
			want:    []int16{1, 2, 5, 6, 5, 6, 7, 8},
			start:   4,
			end:     6,
			shift:   -2,
			succeed: true,
		},
		"should shift right": {
			want:    []int16{1, 2, 3, 4, 3, 4, 7, 8},
			start:   2,
			end:     4,
			shift:   2,
			succeed: true,
		},
		"should not shift above size": {
			want:    []int16{1, 2, 3, 4, 5, 6, 7, 8},
			start:   6,
			end:     8,
			shift:   2,
			succeed: false,
		},
		"should not shift below start": {
			want:    []int16{1, 2, 3, 4, 5, 6, 7, 8},
			start:   0,
			end:     2,
			shift:   -2,
			succeed: false,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			l := New(make([]byte, 16))

			for i, d := range data {
				l.Set(int16(i), d)
			}

			if c.succeed {
				l.ShiftRange(c.start, c.end, c.shift)
			} else {
				should.Panic(t, "should not shift", func() {
					l.ShiftRange(c.start, c.end, c.shift)
				})
			}

			result := make([]int16, 0, len(data))
			for i := range data {
				d := l.Get(int16(i))
				result = append(result, d)
			}

			if !slices.Equal(result, c.want) {
				t.Errorf("incorrect : %v, expected: %v", result, c.want)
			}
		})
	}
}
