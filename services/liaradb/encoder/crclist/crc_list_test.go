package crclist

import (
	"encoding/binary"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/page"
)

type tuple struct {
	a int16
	b int16
	c page.CRC
}

func TestCRCList_Default(t *testing.T) {
	t.Parallel()

	l := New([]byte{})

	if length := l.Length(); length != 0 {
		t.Errorf("incorrect length: %v, expected: %v", length, 0)
	}

	if s := l.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if c := l.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}

func TestCRCList_Push(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 32))

	if length := l.Length(); length != 32 {
		t.Errorf("incorrect length: %v, expected: %v", length, 32)
	}

	if i, ok := l.Push(1, 10, page.RestoreCRC(int32(100))); !ok {
		t.Error("should push")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if s := l.Size(); s != 10 {
		t.Errorf("incorrect size: %v, expected: %v", s, 10)
	}

	if c := l.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	if i, ok := l.Push(2, 20, page.RestoreCRC(int32(200))); !ok {
		t.Error("should push")
	} else if i != 1 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if s := l.Size(); s != 18 {
		t.Errorf("incorrect size: %v, expected: %v", s, 18)
	}

	if c := l.Count(); c != 2 {
		t.Errorf("incorrect count: %v, expected: %v", c, 2)
	}

	if item, ok := l.Item(0); !ok {
		t.Errorf("should have value")
	} else if item.Offset != 1 {
		t.Errorf("incorrect value: %v, expected: %v", item.Offset, 1)
	} else if item.Size != 10 {
		t.Errorf("incorrect value: %v, expected: %v", item.Size, 10)
	} else if item.CRC != page.RestoreCRC(int32(100)) {
		t.Errorf("incorrect value: %v, expected: %v", item.CRC, page.RestoreCRC(int32(100)))
	}

	if item, ok := l.Item(1); !ok {
		t.Errorf("should have value")
	} else if item.Offset != 2 {
		t.Errorf("incorrect value: %v, expected: %v", item.Offset, 2)
	} else if item.Size != 20 {
		t.Errorf("incorrect value: %v, expected: %v", item.Size, 20)
	} else if item.CRC != page.RestoreCRC(int32(200)) {
		t.Errorf("incorrect value: %v, expected: %v", item.CRC, page.RestoreCRC(int32(200)))
	}

	if _, ok := l.Item(2); ok {
		t.Errorf("should not have a value")
	}
}

func TestCRCList_Pop(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 18))

	if _, ok := l.Push(1, 10, page.RestoreCRC(int32(100))); !ok {
		t.Error("should push")
	}

	if _, ok := l.Push(2, 20, page.RestoreCRC(int32(200))); !ok {
		t.Error("should push")
	}

	if item, ok := l.Pop(); !ok {
		t.Error("should pop")
	} else if item.Offset != 2 {
		t.Errorf("incorrect value: %v, expected: %v", item.Offset, 2)
	} else if item.Size != 20 {
		t.Errorf("incorrect value: %v, expected: %v", item.Size, 20)
	} else if item.CRC != page.RestoreCRC(int32(200)) {
		t.Errorf("incorrect value: %v, expected: %v", item.CRC, page.RestoreCRC(int32(200)))
	}

	if s := l.Size(); s != 10 {
		t.Errorf("incorrect size: %v, expected: %v", s, 10)
	}

	if c := l.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	if item, ok := l.Pop(); !ok {
		t.Error("should pop")
	} else if item.Offset != 1 {
		t.Errorf("incorrect value: %v, expected: %v", item.Offset, 1)
	} else if item.Size != 10 {
		t.Errorf("incorrect value: %v, expected: %v", item.Size, 10)
	} else if item.CRC != page.RestoreCRC(int32(100)) {
		t.Errorf("incorrect value: %v, expected: %v", item.CRC, page.RestoreCRC(int32(100)))
	}

	if s := l.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if c := l.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	if _, ok := l.Pop(); ok {
		t.Error("should not pop beyond empty")
	}
}

func TestCRCList_Items(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 42))

	data := []tuple{
		{10, 60, page.RestoreCRC(int32(100))},
		{20, 70, page.RestoreCRC(int32(200))},
		{30, 80, page.RestoreCRC(int32(300))},
		{40, 90, page.RestoreCRC(int32(400))},
		{50, 100, page.RestoreCRC(int32(500))}}

	for _, i := range data {
		if _, ok := l.Push(i.a, i.b, i.c); !ok {
			t.Error("should push")
		}
	}

	result := make([]tuple, 0, len(data))
	for i := range l.Items() {
		result = append(result, tuple{i.Offset, i.Size, i.CRC})
	}

	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

func TestCRCList_ItemsReverse(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 42))

	data := []tuple{
		{10, 60, page.RestoreCRC(int32(100))},
		{20, 70, page.RestoreCRC(int32(200))},
		{30, 80, page.RestoreCRC(int32(300))},
		{40, 90, page.RestoreCRC(int32(400))},
		{50, 100, page.RestoreCRC(int32(500))}}

	for _, i := range data {
		if _, ok := l.Push(i.a, i.b, i.c); !ok {
			t.Error("should push")
		}
	}

	result := make([]tuple, 0, len(data))
	for i := range l.ItemsReverse() {
		result = append(result, tuple{i.Offset, i.Size, i.CRC})
	}

	// Partial iteration
	for range l.ItemsReverse() {
		break
	}

	slices.Reverse(data)
	if !slices.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}

func TestCRCList_ItemsRange(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 42))

	data := []tuple{
		{10, 60, page.RestoreCRC(int32(100))},
		{20, 70, page.RestoreCRC(int32(200))},
		{30, 80, page.RestoreCRC(int32(300))},
		{40, 90, page.RestoreCRC(int32(400))},
		{50, 100, page.RestoreCRC(int32(500))}}

	for _, i := range data {
		if _, ok := l.Push(i.a, i.b, i.c); !ok {
			t.Error("should push")
		}
	}

	for message, c := range map[string]struct {
		skip  bool
		want  []tuple
		start int16
		end   int16
	}{
		"should iterate the range": {
			want: []tuple{
				{20, 70, page.RestoreCRC(int32(200))},
				{30, 80, page.RestoreCRC(int32(300))},
				{40, 90, page.RestoreCRC(int32(400))}},
			start: 1,
			end:   4,
		},
		"should iterate wrapping the end": {
			want:  data,
			start: 0,
			end:   -1,
		},
		"should iterate wrapping the start": {
			want: []tuple{
				{30, 80, page.RestoreCRC(int32(300))},
				{40, 90, page.RestoreCRC(int32(400))}},
			start: -4,
			end:   -2,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			result := make([]tuple, 0, len(c.want))
			for i := range l.ItemsRange(c.start, c.end) {
				result = append(result, tuple{i.Offset, i.Size, i.CRC})
			}

			// Partial iteration
			for range l.ItemsRange(c.start, c.end) {
				break
			}

			if !slices.Equal(result, c.want) {
				t.Errorf("incorrect result: %v, expected: %v", result, c.want)
			}
		})
	}
}

func TestCRCList_Insert(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip   bool
		data   []tuple
		want   []tuple
		insert tuple
		index  int16
	}{
		"should insert into beginning": {
			data: []tuple{
				{20, 70, page.RestoreCRC(int32(200))},
				{30, 80, page.RestoreCRC(int32(300))},
				{40, 90, page.RestoreCRC(int32(400))},
				{50, 100, page.RestoreCRC(int32(500))}},
			want: []tuple{
				{10, 60, page.RestoreCRC(int32(100))},
				{20, 70, page.RestoreCRC(int32(200))},
				{30, 80, page.RestoreCRC(int32(300))},
				{40, 90, page.RestoreCRC(int32(400))},
				{50, 100, page.RestoreCRC(int32(500))}},
			insert: tuple{10, 60, page.RestoreCRC(int32(100))},
			index:  0,
		},
		"should insert into middle": {
			data: []tuple{
				{10, 60, page.RestoreCRC(int32(100))},
				{20, 70, page.RestoreCRC(int32(200))},
				{40, 90, page.RestoreCRC(int32(400))},
				{50, 100, page.RestoreCRC(int32(500))}},
			want: []tuple{
				{10, 60, page.RestoreCRC(int32(100))},
				{20, 70, page.RestoreCRC(int32(200))},
				{30, 80, page.RestoreCRC(int32(300))},
				{40, 90, page.RestoreCRC(int32(400))},
				{50, 100, page.RestoreCRC(int32(500))}},
			insert: tuple{30, 80, page.RestoreCRC(int32(300))},
			index:  2,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			l := New(make([]byte, 42))

			for _, i := range c.data {
				if _, ok := l.Push(i.a, i.b, i.c); !ok {
					t.Fatal("should push")
				}
			}

			if _, ok := l.Insert(c.insert.a, c.insert.b, c.insert.c, c.index); !ok {
				t.Fatal("should insert")
			}

			wantCount := int16(len(c.want))
			if count := l.Count(); count != wantCount {
				t.Errorf("incorrect count: %v, expected: %v", count, wantCount)
			}

			result := make([]tuple, 0, len(c.data))
			for i := range l.Items() {
				result = append(result, tuple{i.Offset, i.Size, i.CRC})
			}

			// Partial iteration
			for range l.Items() {
				break
			}

			if !slices.Equal(result, c.want) {
				t.Errorf("incorrect result: %v, expected: %v", result, c.want)
			}
		})
	}

	t.Run("should not insert beyond size", func(t *testing.T) {
		t.Parallel()

		l := New(make([]byte, 10))
		for i, item := range []tuple{
			{10, 60, page.RestoreCRC(int32(100))}} {
			if _, ok := l.Insert(item.a, item.b, item.c, int16(i)); !ok {
				t.Fatal("should insert")
			}
		}

		if _, ok := l.Insert(20, 70, page.RestoreCRC(int32(200)), 0); ok {
			t.Error("should not insert beyond size")
		}
	})
}

func TestCRCList_Reset(t *testing.T) {
	t.Parallel()

	data := make([]byte, 32)

	l := New(data)

	if s := l.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if c := l.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	binary.BigEndian.PutUint16(data, 3)
	l.Reset()

	if s := l.Size(); s != 26 {
		t.Errorf("incorrect size: %v, expected: %v", s, 26)
	}

	if c := l.Count(); c != 3 {
		t.Errorf("incorrect count: %v, expected: %v", c, 3)
	}

	c := 0
	for range l.Items() {
		c++
	}

	if c != 3 {
		t.Errorf("incorrect count: %v, expected: %v", c, 3)
	}
}

func TestCRCList_Clear(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 16))

	if length := l.Length(); length != 16 {
		t.Errorf("incorrect length: %v, expected: %v", length, 16)
	}

	if i, ok := l.Push(1, 2, page.CRC{}); !ok {
		t.Error("should push")
	} else if i != 0 {
		t.Errorf("incorrect index: %v, expected: %v", i, 0)
	}

	if s := l.Size(); s != 10 {
		t.Errorf("incorrect size: %v, expected: %v", s, 10)
	}

	if c := l.Count(); c != 1 {
		t.Errorf("incorrect count: %v, expected: %v", c, 1)
	}

	l.Clear()

	if s := l.Size(); s != 2 {
		t.Errorf("incorrect size: %v, expected: %v", s, 2)
	}

	if c := l.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	c := 0
	for range l.Items() {
		c++
	}

	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}

func TestCRCList_SetCRC(t *testing.T) {
	t.Parallel()

	l := New(make([]byte, 10))
	if _, ok := l.Insert(10, 60, page.RestoreCRC(int32(100)), 0); !ok {
		t.Fatal("should insert")
	}

	if item, ok := l.Item(0); !ok {
		t.Fatal("should get item")
	} else if item.Offset != 10 {
		t.Errorf("incorrect offset: %v, expected: %v", item.Offset, 10)
	} else if item.Size != 60 {
		t.Errorf("incorrect size: %v, expected: %v", item.Size, 60)
	} else if item.CRC != page.RestoreCRC(int32(100)) {
		t.Errorf("incorrect crc: %v, expected: %v", item.CRC, page.RestoreCRC(int32(100)))
	}

	l.SetCRC(page.RestoreCRC(int32(200)), 0)

	if item, ok := l.Item(0); !ok {
		t.Fatal("should get item")
	} else if item.Offset != 10 {
		t.Errorf("incorrect offset: %v, expected: %v", item.Offset, 10)
	} else if item.Size != 60 {
		t.Errorf("incorrect size: %v, expected: %v", item.Size, 60)
	} else if item.CRC != page.RestoreCRC(int32(200)) {
		t.Errorf("incorrect crc: %v, expected: %v", item.CRC, page.RestoreCRC(int32(200)))
	}
}
