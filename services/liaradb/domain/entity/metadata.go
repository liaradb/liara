package entity

import (
	"io"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
)

// TODO: Map ClientVersion
type Metadata struct {
	UserID        value.UserID        // The ID of the User issuing the Command
	CorrelationID value.CorrelationID // The ID of the entire Command and Event chain
	ClientVersion string              // The Version of the Client emitting the Event
	Time          value.Time          // The Time this Event was created
}

func (e Metadata) Size() int {
	return raw.Size(
		e.UserID,
		e.CorrelationID,
		e.Time)
}

func (e Metadata) Write(w io.Writer) error {
	return raw.WriteAll(w,
		e.UserID,
		e.CorrelationID,
		e.Time)
}

func (e *Metadata) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&e.UserID,
		&e.CorrelationID,
		&e.Time)
}
