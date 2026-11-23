package application

import (
	"github.com/liaradb/liaradb/collection/heap"
	"github.com/liaradb/liaradb/collection/indexed"
	"github.com/liaradb/liaradb/domain/infrastructure/eventlog"
	"github.com/liaradb/liaradb/domain/infrastructure/idempotency"
	"github.com/liaradb/liaradb/domain/infrastructure/keyvalue"
	"github.com/liaradb/liaradb/domain/infrastructure/manager"
	"github.com/liaradb/liaradb/domain/infrastructure/outbox"
	"github.com/liaradb/liaradb/domain/infrastructure/tenant"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/transaction"
)

type services struct {
	es *service.EventService
	os *service.OutboxService
	ts *service.TenantService
}

type collections struct {
	manager     *manager.Manager
	Tenant      *tenant.Tenant
	EventLog    *eventlog.EventLog
	KeyValue    *keyvalue.KeyValue
	Outbox      *outbox.Outbox
	Idempotency *idempotency.Idempotency
}

func newServices(txManager *transaction.Manager) *services {
	c := newCollections()
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

func newCollections() *collections {
	ic := indexed.New(heap.New())
	kv := keyvalue.New(ic)
	return &collections{
		manager:     manager.New(kv),
		Tenant:      tenant.New(ic),
		EventLog:    eventlog.New(ic),
		KeyValue:    kv,
		Outbox:      outbox.New(ic),
		Idempotency: idempotency.New(ic),
	}
}
