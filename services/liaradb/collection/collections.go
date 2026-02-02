package collection

import (
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/idempotency"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/manager"
	"github.com/liaradb/liaradb/collection/outbox"
	"github.com/liaradb/liaradb/collection/tenant"
	"github.com/liaradb/liaradb/storage"
)

type Collections struct {
	storage     *storage.Storage
	manager     *manager.Manager
	tenant      *tenant.Tenant
	EventLog    *eventlog.EventLog
	keyValue    *keyvalue.KeyValue
	outbox      *outbox.Outbox
	idempotency *idempotency.Idempotency
}

func NewCollections(
	storage *storage.Storage,
) *Collections {
	cursor := btree.NewCursor(storage)
	kv := keyvalue.New(storage, cursor)
	return &Collections{
		storage:     storage,
		manager:     manager.New(kv),
		tenant:      tenant.New(storage, cursor),
		EventLog:    eventlog.New(storage, cursor),
		keyValue:    kv,
		outbox:      outbox.New(storage, cursor),
		idempotency: idempotency.New(storage, cursor),
	}
}
