package command

import (
	"time"

	"github.com/liaradb/liaradb/domain/value"
)

type AppendEvent struct {
	TenantID      value.TenantID
	RequestID     *value.RequestID    // The ID of the Request, for idempotency
	CorrelationID value.CorrelationID // The ID of the entire Command and Event chain
	ClientVersion value.ClientVersion // The Version of the client
	UserID        value.UserID        // The ID of the User issuing the Command
	Time          time.Time           // The Time this Event was created
	Events        []EventOptions
}
