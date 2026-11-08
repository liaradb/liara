package application

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
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction"
)

type services struct {
	es *service.EventService
	os *service.OutboxService
	ts *service.TenantService
}

type collections struct {
	manager     *manager.Manager
	schemaMgr   *schema.Manager
	Tenant      *tenant.Tenant
	EventLog    *eventlog.EventLog
	KeyValue    *keyvalue.KeyValue
	Outbox      *outbox.Outbox
	Idempotency *idempotency.Idempotency
}

func newServices(txManager *transaction.Manager, s *storage.Storage) *services {
	c := newCollections(s)
	return &services{
		es: service.NewEventService(
			txManager,
			c.Tenant,
			c.Idempotency,
			c.EventLog,
			c.KeyValue,
			c.Outbox),
		os: service.NewOutboxService(
			txManager,
			c.Tenant,
			c.Idempotency,
			c.Outbox),
		ts: service.NewTenantService(
			txManager,
			c.Tenant),
	}
}

func newCollections(s *storage.Storage) *collections {
	c := btree.NewCursor(s)
	fc := fixed.New(record.New(s), c)
	kv := keyvalue.New(fc)
	return &collections{
		manager:     manager.New(kv),
		schemaMgr:   schema.NewManager(s),
		Tenant:      tenant.New(fc),
		EventLog:    eventlog.New(fc, c),
		KeyValue:    kv,
		Outbox:      outbox.New(fc),
		Idempotency: idempotency.New(fc, c),
	}
}
