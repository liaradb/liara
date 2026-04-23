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

	o, err := esc.GetOutbox(t.Context(), &liara.GetOutboxRequest{
		OutboxId: res0.OutboxId,
		TenantId: tid,
	})
	if err != nil {
		t.Fatal(err)
	}

	if o.GlobalVersion != 0 {
		t.Errorf("incorrect global version: %v, expected: %v", o.GlobalVersion, 0)
	}

	listReq := &liara.ListOutboxesRequest{
		TenantId: tid,
	}

	var result []*liara.Outbox
	if err := esc.ListOutboxes(listReq, newTestStream(t.Context(), func(o *liara.Outbox) {
		result = append(result, o)
	})); err != nil {
		t.Error(err)
	}

	if len(result) != 1 {
		t.Errorf("incorrect length: %v, expected: %v", len(result), 1)
	}

	o, err = esc.GetOutbox(t.Context(), &liara.GetOutboxRequest{
		OutboxId: res0.OutboxId,
		TenantId: tid,
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := esc.UpdateOutboxPosition(t.Context(), &liara.UpdateOutboxPositionRequest{
		TenantId:      tid,
		OutboxId:      res0.OutboxId,
		GlobalVersion: 10,
	}); err != nil {
		t.Fatal(err)
	}

	if o.GlobalVersion != 0 {
		t.Errorf("incorrect global version: %v, expected: %v", o.GlobalVersion, 10)
	}
}

func TestEventSourceController__Event(t *testing.T) {
	esc := NewEventSourceController(&testEventService{}, &testTenantService{})

	tid := uuid.NewString()
	var pid int32 = 1
	id := uuid.NewString()
	name := "name"

	if _, err := esc.Append(t.Context(), &liara.AppendRequest{
		TenantId:    tid,
		PartitionId: pid,
		Events: []*liara.AppendEvent{{
			AggregateId: id,
			Id:          uuid.NewString(),
			Version:     1,
		}},
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := esc.Append(t.Context(), &liara.AppendRequest{
		TenantId:    tid,
		PartitionId: pid,
		Events: []*liara.AppendEvent{{
			AggregateId:   id,
			AggregateName: name,
			Id:            uuid.NewString(),
			Version:       2,
		}},
	}); err != nil {
		t.Fatal(err)
	}

	var result []*liara.Event
	if err := esc.Get(&liara.GetRequest{
		TenantId:    tid,
		PartitionId: pid,
		AggregateId: id,
	}, newTestStream(t.Context(), func(e *liara.Event) {
		result = append(result, e)
	})); err != nil {
		t.Fatal(err)
	}

	if l := len(result); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	var result2 []*liara.Event
	if err := esc.GetByAggregateIDAndName(&liara.GetByAggregateIDAndNameRequest{
		TenantId:    tid,
		PartitionId: pid,
		AggregateId: id,
		Name:        name,
	}, newTestStream(t.Context(), func(e *liara.Event) {
		result2 = append(result2, e)
	})); err != nil {
		t.Fatal(err)
	}

	if l := len(result2); l != 1 {
		t.Errorf("incorrect length: %v, expected: %v", l, 1)
	}
}
