package collection

import (
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/idempotency"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/manager"
	"github.com/liaradb/liaradb/collection/outbox"
	"github.com/liaradb/liaradb/collection/schema"
	"github.com/liaradb/liaradb/collection/tenant"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction/log"
)

type Collections struct {
	storage     *storage.Storage
	manager     *manager.Manager
	schemaMgr   *schema.Manager
	tenant      *tenant.Tenant
	EventLog    *eventlog.EventLog
	keyValue    *keyvalue.KeyValue
	Outbox      *outbox.Outbox
	Idempotency *idempotency.Idempotency
}

func NewCollections(
	storage *storage.Storage,
	log *log.Log,
) *Collections {
	cursor := btree.NewCursor(storage)
	kv := keyvalue.New(storage, cursor, log)
	return &Collections{
		storage:     storage,
		manager:     manager.New(kv),
		tenant:      tenant.New(storage, cursor, log),
		EventLog:    eventlog.New(storage, cursor, log),
		keyValue:    kv,
		Outbox:      outbox.New(storage, cursor, log),
		Idempotency: idempotency.New(storage, cursor, log),
	}
}
