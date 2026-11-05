package application

import (
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/transaction"
)

type services struct {
	es *service.EventService
	os *service.OutboxService
	ts *service.TenantService
}

func newServices(txManager *transaction.Manager, c *collections) *services {
	return &services{
		es: service.NewEventService(
			txManager,
			c.Idempotency,
			c.EventLog,
			c.KeyValue,
			c.Outbox),
		os: service.NewOutboxService(
			txManager,
			c.Idempotency,
			c.Outbox),
		ts: service.NewTenantService(
			txManager,
			c.Tenant),
	}
}
