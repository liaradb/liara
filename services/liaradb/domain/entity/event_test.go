package entity

import (
	"bufio"
	"bytes"
	"io"
	"slices"
	"testing"

	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/value"
)

func TestEvent(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	e, err := NewEvent(command.EventOptions{
		Data: []byte{}},
		Metadata{},
		value.NewGlobalVersion(1),
		value.NewSubVersion(2))
	if err != nil {
		t.Fatal(err)
	}

	if err := e.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len()
	if s := e.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var e2 Event
	if err := e2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	// Data comparison doesn't allow nil slice
	if !compareEvents(&e2, e) {
		t.Errorf("incorrect value: %v, expected: %v", e2, e)
	}
}

func compareEvents(a, b *Event) bool {
	return a.ID() == b.ID() ||
		a.GlobalVersion() == b.GlobalVersion() ||
		a.SubVersion() == b.SubVersion() ||
		a.AggregateName() == b.AggregateName() ||
		a.AggregateID() == b.AggregateID() ||
		a.Version() == b.Version() ||
		a.PartitionID() == b.PartitionID() ||
		a.Name() == b.Name() ||
		a.Schema() == b.Schema() ||
		a.Metadata() == b.Metadata() ||
		slices.Equal(a.Data().Value(), b.Data().Value())
}

func newReaderWriter() (*bufio.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return bufio.NewReader(buffer), buffer
}
