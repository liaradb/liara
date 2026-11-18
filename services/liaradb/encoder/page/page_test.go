package page

import (
	"fmt"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/encoder/scan"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/testing/should"
)

func TestPage_New(t *testing.T) {
	t.Parallel()

	p := New(32, 4, 4)

	want := make([]byte, 32)
	_ = scan.SetInt32(want, int32(MagicPage))
	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_NewFromSlice(t *testing.T) {
	t.Parallel()

	want := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 0, 0, 0, 0}
	p := NewFromSlice(want, 4, 4)

	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_Fill(t *testing.T) {
	t.Parallel()

	p := New(20, 4, 4)

	want := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 0, 0, 0, 0}
	p.Fill(want)

	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}

	p.Clear()

	want = make([]byte, 20)
	_ = scan.SetInt32(want, int32(MagicPage))
	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_Replace(t *testing.T) {
	t.Parallel()

	p := New(20, 4, 4)

	want := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 0, 0, 0, 0}
	b := buffer.NewFromSlice(want)
	if err := p.Replace(b); err != nil {
		t.Fatal(err)
	}

	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}

	p.Clear()

	want = make([]byte, 20)
	_ = scan.SetInt32(want, int32(MagicPage))
	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_Replace__Error(t *testing.T) {
	t.Parallel()

	p := New(20, 4, 4)

	if err := p.Replace(&errorReader{}); err == nil {
		t.Error("should return error")
	}
}

type errorReader struct {
}

func (r *errorReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("error")
}

func TestPage_Header(t *testing.T) {
	t.Parallel()

	want := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 0, 0, 0, 0}
	p := NewFromSlice(want, 4, 4)

	if h := p.Header(); !slices.Equal(h, want[headerSize:headerSize+4]) {
		t.Errorf("incorrect data: %v, expected: %v", h, want[headerSize:headerSize+4])
	}

	if tid := p.TrackingID(); tid != 0 {
		t.Errorf("incorrect tracking id: %v, expected: %v", tid, 0)
	}

	p.SetTrackingID(1)

	if tid := p.TrackingID(); tid != 1 {
		t.Errorf("incorrect tracking id: %v, expected: %v", tid, 1)
	}
}

func TestPage_Next(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)

	header, data := p.Next(8)

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	_ = p.CommitAppend(8)

	header, data = p.Slot(0)

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}
}

func TestPage_Next__Empty(t *testing.T) {
	t.Parallel()

	p := New(32, 4, 4)
	_ = p.CommitAppend(18)

	header, data := p.Next(2)
	if l := len(header); l != 0 {
		t.Errorf("incorrect header: %v, expected: %v", l, 0)
	}

	if l := len(data); l != 0 {
		t.Errorf("incorrect data: %v, expected: %v", l, 0)
	}
}

func TestPage_NextMustFit(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)

	if _, _, ok := p.NextMustFit(100); ok {
		t.Error("should not fit oversized")
	}

	header, data, _ := p.NextMustFit(8)

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}

	_ = p.CommitAppend(8)

	header, data = p.Slot(0)

	if l := len(header); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if l := len(data); l != 8 {
		t.Errorf("incorrect length: %v, expected: %v", l, 8)
	}
}

func TestPage_Insert(t *testing.T) {
	t.Parallel()

	p := New(128, 4, 4)

	if c := p.Count(); c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	{
		header, data := p.Next(8)

		if l := len(header); l != 4 {
			t.Errorf("incorrect length: %v, expected: %v", l, 4)
		}

		if l := len(data); l != 8 {
			t.Errorf("incorrect length: %v, expected: %v", l, 8)
		}

		header[0] = 1
		data[0] = 2

		_ = p.CommitAppend(8)

		if c := p.Count(); c != 1 {
			t.Errorf("incorrect count: %v, expected: %v", c, 1)
		}
	}

	{
		header, data := p.Next(8)

		if l := len(header); l != 4 {
			t.Errorf("incorrect length: %v, expected: %v", l, 4)
		}

		if l := len(data); l != 8 {
			t.Errorf("incorrect length: %v, expected: %v", l, 8)
		}

		header[0] = 3
		data[0] = 4

		_ = p.Insert(8, 0)

		if c := p.Count(); c != 2 {
			t.Errorf("incorrect count: %v, expected: %v", c, 2)
		}
	}

	{
		header, data := p.Slot(0)

		if l := len(header); l != 4 {
			t.Errorf("incorrect length: %v, expected: %v", l, 4)
		}

		if l := len(data); l != 8 {
			t.Errorf("incorrect length: %v, expected: %v", l, 8)
		}

		if h := header[0]; h != 3 {
			t.Errorf("incorrect header: %v, expected: %v", h, 3)
		}

		if d := data[0]; d != 4 {
			t.Errorf("incorrect data: %v, expected: %v", d, 4)
		}
	}

	{
		header, data := p.Slot(1)

		if l := len(header); l != 4 {
			t.Errorf("incorrect length: %v, expected: %v", l, 4)
		}

		if l := len(data); l != 8 {
			t.Errorf("incorrect length: %v, expected: %v", l, 8)
		}

		if h := header[0]; h != 1 {
			t.Errorf("incorrect header: %v, expected: %v", h, 1)
		}

		if d := data[0]; d != 2 {
			t.Errorf("incorrect data: %v, expected: %v", d, 2)
		}
	}
}

func TestPage_Slot(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)

	should.Panic(t, "slot should not exist", func() {
		_, _ = p.Slot(0)
	})

	c := 0
	for range p.Slots() {
		c++
	}

	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}

	_, _ = p.Next(8)
	_ = p.CommitAppend(8)

	_, _ = p.Next(8)
	_ = p.CommitAppend(8)

	_, _ = p.Slot(0)
	_, _ = p.Slot(1)

	c = 0
	for range p.Slots() {
		c++
	}

	if c != 2 {
		t.Errorf("incorrect count: %v, expected: %v", c, 2)
	}

	// Early return
	for range p.Slots() {
		break
	}
}

func TestPage_Slots(t *testing.T) {
	t.Parallel()

	p := New(256, 4, 4)

	want := []int64{101, 102, 103, 104, 105}
	for _, v := range want {
		_, b := p.Next(8)
		_ = scan.SetInt64(b, v)
		_ = p.CommitAppend(8)
	}

	result := make([]int64, 0, 5)
	for _, b := range p.Slots() {
		v, _ := scan.Int64(b)
		result = append(result, v)
	}

	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}
}

func TestPage_SlotsRange(t *testing.T) {
	t.Parallel()

	p := New(256, 4, 4)

	want := []int64{101, 102, 103, 104, 105}
	for _, v := range want {
		_, b := p.Next(8)
		_ = scan.SetInt64(b, v)
		_ = p.CommitAppend(8)
	}

	result := make([]int64, 0, 5)
	for _, b := range p.SlotsRange(1, 4) {
		v, _ := scan.Int64(b)
		result = append(result, v)
	}

	if !slices.Equal(result, want[1:4]) {
		t.Errorf("incorrect result: %v, expected: %v", result, want[1:4])
	}

	for range p.SlotsRange(1, 4) {
		// Early return
		break
	}
}

func TestPage_SlotsReverse(t *testing.T) {
	t.Parallel()

	p := New(256, 4, 4)

	want := []int64{101, 102, 103, 104, 105}
	for _, v := range want {
		_, b := p.Next(8)
		_ = scan.SetInt64(b, v)
		_ = p.CommitAppend(8)
	}

	result := make([]int64, 0, 5)
	for _, b := range p.SlotsReverse() {
		v, _ := scan.Int64(b)
		result = append(result, v)
	}

	slices.Reverse(want)
	if !slices.Equal(result, want) {
		t.Errorf("incorrect result: %v, expected: %v", result, want)
	}

	for range p.SlotsReverse() {
		// Early return
		break
	}
}

func TestPage_Delete(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)

	_, _ = p.Next(8)
	_ = p.CommitAppend(8)

	_, _ = p.Slot(0)

	p.Delete(0)
	should.Panic(t, "slot should not exist", func() {
		_, _ = p.Slot(0)
	})
}

func TestPage_RecoverAppend(t *testing.T) {
	t.Parallel()

	p0 := New(40, 4, 4)

	header, data := p0.Next(8)
	header[1] = 1
	data[2] = 2

	_ = p0.CommitAppend(8)

	slots0 := make([][]byte, 0)
	for h, b := range p0.Slots() {
		slot := append(slices.Clone(h), b...)
		slots0 = append(slots0, slot)
	}

	p1 := New(40, 4, 4)
	for _, slot := range slots0 {
		p1.RecoverAppend(slot)
	}

	slots1 := make([][]byte, 0)
	for h, b := range p1.Slots() {
		slot := append(slices.Clone(h), b...)
		slots1 = append(slots1, slot)
	}

	if !slices.EqualFunc(slots1, slots0, slices.Equal) {
		t.Errorf("incorrect data:\n%v,\nexpected:\n%v", slots0, slots1)
	}
}

func TestPage_RevertAppend(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)

	base := slices.Clone(p.Data())

	header, data := p.Next(8)
	header[1] = 1
	data[2] = 2

	sid := p.CommitAppend(8)
	h, b := p.Slot(sid)
	slot := append(slices.Clone(h), b...)

	version0 := slices.Clone(p.Data())
	p.RevertAppend()
	if !slices.Equal(p.Data(), base) {
		t.Errorf("incorrect base: %v, expected: %v", p.Data(), base)
	}

	p.RecoverAppend(slot)
	if !slices.Equal(p.Data(), version0) {
		t.Errorf("incorrect data:\n%v,\nexpected:\n%v", p.Data(), version0)
	}
}

func TestPage_RevertAppend__Panic(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)

	should.Panic(t, "should panic", func() {
		p.RecoverAppend(make([]byte, 50))
	})
}

func TestPage_RecoverInsert(t *testing.T) {
	t.Parallel()

	p0 := New(1024, 4, 4)

	{
		header, data := p0.Next(8)
		header[1] = 1
		data[2] = 2

		_ = p0.Insert(8, link.SlotID(0))
	}
	{
		header, data := p0.Next(8)
		header[1] = 3
		data[2] = 4

		_ = p0.Insert(8, link.SlotID(0))
	}

	slots0 := make([][]byte, 0)
	for h, b := range p0.Slots() {
		slot := append(slices.Clone(h), b...)
		slots0 = append(slots0, slot)
	}

	reverse := slices.Clone(slots0)
	slices.Reverse(reverse)

	p1 := New(1024, 4, 4)
	for _, slot := range reverse {
		p1.RecoverInsert(link.SlotID(0), slot)
	}

	slots1 := make([][]byte, 0)
	for h, b := range p1.Slots() {
		slot := append(slices.Clone(h), b...)
		slots1 = append(slots1, slot)
	}

	if !slices.EqualFunc(slots1, slots0, slices.Equal) {
		t.Errorf("incorrect data:\n%v,\nexpected:\n%v", slots0, slots1)
	}
}

func TestPage_RecoverInsert__Panic(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)

	should.Panic(t, "should panic", func() {
		p.RecoverInsert(0, make([]byte, 50))
	})
}

func TestPage_RecoverUpdate__Panic(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)

	p.CommitAppend(8)

	should.Panic(t, "should panic", func() {
		p.RecoverUpdate(0, make([]byte, 50))
	})
}

func TestPage_RecoverHeader(t *testing.T) {
	t.Parallel()

	p0 := New(40, 4, 4)

	header0 := p0.Header()
	header0[1] = 1

	p1 := New(40, 4, 4)
	p1.RecoverHeader(header0)

	header1 := p1.Header()

	if !slices.Equal(header1, header0) {
		t.Errorf("incorrect data:\n%v,\nexpected:\n%v", header1, header0)
	}
}

func TestPage_RecoverHeader__Panic(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)

	should.Panic(t, "should panic", func() {
		p.RecoverHeader(make([]byte, 50))
	})
}

func TestPage_RecoverClear(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)
	want := slices.Clone(p.Data())

	p.CommitAppend(8)

	p.RecoverClear()

	if data := p.Data(); !slices.Equal(data, want) {
		t.Errorf("incorrect data: %v, expected: %v", data, want)
	}
}

func TestPage_IsPage(t *testing.T) {
	t.Parallel()

	p := New(40, 4, 4)

	if !p.IsPage() {
		t.Error("should be page")
	}
	if p.IsEmpty() {
		t.Error("should not be empty")
	}

	p.Fill(make([]byte, 40))

	if p.IsPage() {
		t.Error("should not be page")
	}
	if !p.IsEmpty() {
		t.Error("should be empty")
	}
}
