package controller

import (
	"testing"

	"github.com/google/uuid"
	liara "github.com/liaradb/eventsource_go/generated"
)

func TestEventSourceController__Outbox(t *testing.T) {
	esc := NewEventSourceController(&testEventService{}, &testTenantService{})

	tid := uuid.NewString()

	res0, err := esc.CreateOutbox(t.Context(), &liara.CreateOutboxRequest{
		TenantId: tid,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = esc.GetOutbox(t.Context(), &liara.GetOutboxRequest{
		OutboxId: res0.OutboxId,
		TenantId: tid,
	})
	if err != nil {
		t.Fatal(err)
	}
}
