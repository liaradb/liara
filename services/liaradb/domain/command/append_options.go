package command

import (
	"time"

	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type AppendOptions struct {
	requestID     *value.RequestID    // The ID of the Request, for idempotency
	correlationID value.CorrelationID // The ID of the entire Command and Event chain
	clientVersion value.ClientVersion // The Version of the client
	userID        value.UserID        // The ID of the User issuing the Command
	time          time.Time           // The Time this Event was created
}

func NewAppendOptions(
	requestID *value.RequestID, // The ID of the Request, for idempotency
	correlationID value.CorrelationID, // The ID of the entire Command and Event chain
	clientVersion value.ClientVersion, // The Version of the client
	userID value.UserID, // The ID of the User issuing the Command
	t time.Time, // The Time this Event was created
) AppendOptions {
	ao := AppendOptions{
		requestID:     requestID,
		correlationID: correlationID,
		clientVersion: clientVersion,
		userID:        userID,
		time:          t,
	}

	if ao.Time().IsZero() {
		ao.time = time.Now()
	}

	return ao
}

func (ao *AppendOptions) RequestID() (value.RequestID, bool) {
	if ao.requestID == nil {
		return value.RequestID{}, false
	}

	return *ao.requestID, true
}

func (ao *AppendOptions) CorrelationID() value.CorrelationID { return ao.correlationID }
func (ao *AppendOptions) ClientVersion() value.ClientVersion { return ao.clientVersion }
func (ao *AppendOptions) UserID() value.UserID               { return ao.userID }
func (ao *AppendOptions) Time() time.Time                    { return ao.time }

func (ao *AppendOptions) toMetadata() entity.Metadata {
	return entity.NewMetadata(
		ao.userID,
		ao.correlationID,
		ao.clientVersion,
		value.NewTime(ao.time),
	)
}
