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

	listReq := &liara.ListOutboxesRequest{
		TenantId: tid,
	}

	var result []*liara.Outbox
	listStr := newTestOutboxStream(t.Context(), func(o *liara.Outbox) {
		result = append(result, o)
	})

	if err := esc.ListOutboxes(listReq, listStr); err != nil {
		t.Error(err)
	}

	if len(result) != 1 {
		t.Errorf("incorrect length: %v, expected: %v", len(result), 1)
	}
}
