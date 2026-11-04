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
)

type Collections struct {
	storage     *storage.Storage
	manager     *manager.Manager
	schemaMgr   *schema.Manager
	Tenant      *tenant.Tenant
	EventLog    *eventlog.EventLog
	KeyValue    *keyvalue.KeyValue
	Outbox      *outbox.Outbox
	Idempotency *idempotency.Idempotency
}

func NewCollections(s *storage.Storage) *Collections {
	c := btree.NewCursor(s)
	kv := keyvalue.New(s, c)
	return &Collections{
		storage:     s,
		manager:     manager.New(kv),
		Tenant:      tenant.New(s, c),
		EventLog:    eventlog.New(s, c),
		KeyValue:    kv,
		Outbox:      outbox.New(s, c),
		Idempotency: idempotency.New(s, c),
	}
}
