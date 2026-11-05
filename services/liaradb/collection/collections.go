package collection

import (
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/collection/fixed"
	"github.com/liaradb/liaradb/collection/idempotency"
	"github.com/liaradb/liaradb/collection/keyvalue"
	"github.com/liaradb/liaradb/collection/manager"
	"github.com/liaradb/liaradb/collection/outbox"
	"github.com/liaradb/liaradb/collection/record"
	"github.com/liaradb/liaradb/collection/schema"
	"github.com/liaradb/liaradb/collection/tenant"
	"github.com/liaradb/liaradb/storage"
)

type Collections struct {
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
	fc := fixed.New(record.New(s), c)
	kv := keyvalue.New(fc)
	return &Collections{
		manager:     manager.New(kv),
		schemaMgr:   schema.NewManager(s),
		Tenant:      tenant.New(fc),
		EventLog:    eventlog.New(fc, c),
		KeyValue:    kv,
		Outbox:      outbox.New(fc),
		Idempotency: idempotency.New(fc),
	}
}
