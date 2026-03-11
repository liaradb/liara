package entity

import (
	"io"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/serializer"
)

type Metadata struct {
	userID        value.UserID        // The ID of the User issuing the Command
	correlationID value.CorrelationID // The ID of the entire Command and Event chain
	clientVersion value.ClientVersion // The Version of the Client emitting the Event
	time          value.Time          // The Time this Event was created
}

func NewMetadata(
	uid value.UserID,
	cid value.CorrelationID,
	cv value.ClientVersion,
	t value.Time,
) Metadata {
	return Metadata{
		userID:        uid,
		correlationID: cid,
		clientVersion: cv,
		time:          t,
	}
}

func (e Metadata) UserID() value.UserID               { return e.userID }
func (e Metadata) CorrelationID() value.CorrelationID { return e.correlationID }
func (e Metadata) ClientVersion() value.ClientVersion { return e.clientVersion }
func (e Metadata) Time() value.Time                   { return e.time }

func (e Metadata) Size() int {
	return serializer.Size(
		e.userID,
		e.correlationID,
		e.clientVersion,
		e.time)
}

func (e Metadata) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		e.userID,
		e.correlationID,
		e.clientVersion,
		e.time)
}

func (e *Metadata) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&e.userID,
		&e.correlationID,
		&e.clientVersion,
		&e.time)
}
