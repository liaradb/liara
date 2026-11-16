package entity

import (
	"bufio"
	"bytes"
	"io"
	"slices"
	"testing"
	"uuid"

	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/value"
)

func TestEvent(t *testing.T) {
	t.Parallel()

	var (
		eventID       = uuid.NewV7().String()
		aggregateName = value.NewAggregateName("aggregate")
		aggregateID   = value.NewAggregateID("aggregate_id")
		version       = value.NewVersion(10)
		partitionID   = value.NewPartitionID(101)
		name          = value.NewEventName("event")
		schema        = value.NewSchema("0.0.0")
		metadata      = Metadata{
			userID: value.NewUserID("user_id"),
		}
		data = []byte{1, 2, 3}
	)

	e, err := NewEvent(command.EventOptions{
		ID:            eventID,
		AggregateName: aggregateName,
		AggregateID:   aggregateID,
		Version:       version,
		PartitionID:   partitionID,
		Name:          name,
		Schema:        schema,
		Data:          data,
	},
		metadata,
		value.NewGlobalVersion(1),
		value.NewSubVersion(2))
	if err != nil {
		t.Fatal(err)
	}

	if i := e.ID().String(); i != eventID {
		t.Errorf("incorrect event id: %v, expected: %v", i, eventID)
	}

	if v := e.GlobalVersion().Value(); v != 1 {
		t.Errorf("incorrect global version: %v, expected: %v", v, 1)
	}

	if v := e.SubVersion().Value(); v != 2 {
		t.Errorf("incorrect sub version: %v, expected: %v", v, 2)
	}

	if n := e.AggregateName(); n != aggregateName {
		t.Errorf("incorrect aggregate name: %v, expected: %v", n, aggregateName)
	}

	if i := e.AggregateID(); i != aggregateID {
		t.Errorf("incorrect aggregate id: %v, expected: %v", i, aggregateName)
	}

	if v := e.Version(); v != version {
		t.Errorf("incorrect version: %v, expected: %v", v, version)
	}

	if i := e.PartitionID(); i != partitionID {
		t.Errorf("incorrect partition id: %v, expected: %v", i, partitionID)
	}

	if n := e.Name(); n != name {
		t.Errorf("incorrect event name: %v, expected: %v", n, name)
	}

	if s := e.Schema(); s != schema {
		t.Errorf("incorrect schema: %v, expected: %v", s, schema)
	}

	if d := e.Data().Value(); !slices.Equal(d, data) {
		t.Errorf("incorrect schema: %v, expected: %v", d, data)
	}

	if i := e.Metadata().UserID(); i != metadata.UserID() {
		t.Errorf("incorrect user id: %v, expected: %v", i, metadata.UserID())
	}
}

func TestEvent_ID(t *testing.T) {
	t.Run("should return error if id is not UUID", func(t *testing.T) {
		_, err := NewEvent(command.EventOptions{
			ID: "random string"},
			Metadata{},
			value.NewGlobalVersion(1),
			value.NewSubVersion(2))
		if err == nil {
			t.Fatal("should return error")
		}
	})

	t.Run("should create event id if empty", func(t *testing.T) {
		e, err := NewEvent(command.EventOptions{},
			Metadata{},
			value.NewGlobalVersion(1),
			value.NewSubVersion(2))
		if err != nil {
			t.Fatal(err)
		}

		if i := e.ID().String(); i == uuid.Nil().String() {
			t.Errorf("incorrect event id: %v, expected: %v", i, uuid.Nil().String())
		}
	})

	t.Run("should use event id if not empty", func(t *testing.T) {
		var (
			eventID = uuid.NewV7().String()
		)

		e, err := NewEvent(command.EventOptions{
			ID: eventID},
			Metadata{},
			value.NewGlobalVersion(1),
			value.NewSubVersion(2))
		if err != nil {
			t.Fatal(err)
		}

		if i := e.ID().String(); i != eventID {
			t.Errorf("incorrect event id: %v, expected: %v", i, eventID)
		}
	})
}

func TestEvent_ReadWrite(t *testing.T) {
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
