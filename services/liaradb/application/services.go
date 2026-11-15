package application

import (
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/heap"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/collection/table/schema"
	"github.com/liaradb/liaradb/domain/infrastructure/eventlog"
	"github.com/liaradb/liaradb/domain/infrastructure/idempotency"
	"github.com/liaradb/liaradb/domain/infrastructure/keyvalue"
	"github.com/liaradb/liaradb/domain/infrastructure/manager"
	"github.com/liaradb/liaradb/domain/infrastructure/outbox"
	"github.com/liaradb/liaradb/domain/infrastructure/tenant"
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
	ic := indexed.New(heap.New(s), c)
	kv := keyvalue.New(ic)
	return &collections{
		manager:     manager.New(kv),
		schemaMgr:   schema.NewManager(s),
		Tenant:      tenant.New(ic),
		EventLog:    eventlog.New(ic, c),
		KeyValue:    kv,
		Outbox:      outbox.New(ic),
		Idempotency: idempotency.New(ic, c),
	}
}
