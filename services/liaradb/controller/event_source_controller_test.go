package controller

import (
	"testing"

	"github.com/google/uuid"
	liara "github.com/liaradb/eventsource_go/generated"
)

func TestEventSourceController__Outbox(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
		Options: &liara.AppendOptions{},
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

func TestEventSourceController__Idempotency(t *testing.T) {
	t.Parallel()

	esc := NewEventSourceController(&testEventService{}, &testTenantService{})

	tid := uuid.NewString()
	var pid int32 = 1
	id := uuid.NewString()
	rid0 := uuid.NewString()

	if _, err := esc.Append(t.Context(), &liara.AppendRequest{
		TenantId:    tid,
		PartitionId: pid,
		Events: []*liara.AppendEvent{{
			AggregateId: id,
			Id:          uuid.NewString(),
			Version:     1,
		}},
		Options: &liara.AppendOptions{
			RequestId: rid0,
		},
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := esc.Append(t.Context(), &liara.AppendRequest{
		TenantId:    tid,
		PartitionId: pid,
		Events: []*liara.AppendEvent{{
			AggregateId: id,
			Id:          uuid.NewString(),
			Version:     2,
		}},
		Options: &liara.AppendOptions{
			RequestId: rid0,
		},
	}); err == nil {
		t.Error("should return error")
	}

	res0, err := esc.TestIdempotency(t.Context(), &liara.TestIdempotencyRequest{
		TenantId:  tid,
		RequestId: rid0,
	})
	if err != nil {
		t.Error(err)
	}
	if res0.Ok {
		t.Errorf("incorrect result: %v, expected: %v", res0.Ok, false)
	}

	res1, err := esc.TestIdempotency(t.Context(), &liara.TestIdempotencyRequest{
		TenantId:  tid,
		RequestId: uuid.NewString(),
	})
	if err != nil {
		t.Error(err)
	}
	if !res1.Ok {
		t.Errorf("incorrect result: %v, expected: %v", res1.Ok, true)
	}

	if _, err = esc.TestIdempotency(t.Context(), &liara.TestIdempotencyRequest{
		TenantId:  "abcde",
		RequestId: uuid.NewString(),
	}); err == nil {
		t.Error("should return error")
	}

	if _, err = esc.TestIdempotency(t.Context(), &liara.TestIdempotencyRequest{
		TenantId:  uuid.NewString(),
		RequestId: "abcde",
	}); err == nil {
		t.Error("should return error")
	}
}

func TestEventSourceController__Tenant(t *testing.T) {
	t.Parallel()

	esc := NewEventSourceController(&testEventService{}, &testTenantService{})

	name0 := "name0"
	name1 := "name1"

	_, err := esc.GetTenant(t.Context(), &liara.GetTenantRequest{
		TenantId: "abcde",
	})
	if err == nil {
		t.Error("should return error")
	}

	_, err = esc.GetTenant(t.Context(), &liara.GetTenantRequest{
		TenantId: uuid.NewString(),
	})
	if err == nil {
		t.Error("should return error")
	}

	_, err = esc.RenameTenant(t.Context(), &liara.RenameTenantRequest{
		TenantId: "abcde",
		Name:     name1,
	})
	if err == nil {
		t.Error("should return error")
	}

	_, err = esc.RenameTenant(t.Context(), &liara.RenameTenantRequest{
		TenantId: uuid.NewString(),
		Name:     name1,
	})
	if err == nil {
		t.Error("should return error")
	}

	_, err = esc.DeleteTenant(t.Context(), &liara.DeleteTenantRequest{
		TenantId: "abcde",
	})
	if err == nil {
		t.Error("should return error")
	}

	_, err = esc.DeleteTenant(t.Context(), &liara.DeleteTenantRequest{
		TenantId: uuid.NewString(),
	})
	if err == nil {
		t.Error("should return error")
	}

	res0, err := esc.CreateTenant(t.Context(), &liara.CreateTenantRequest{
		Name: name0,
	})
	if err != nil {
		t.Error(err)
	}

	res1, err := esc.GetTenant(t.Context(), &liara.GetTenantRequest{
		TenantId: res0.TenantId,
	})
	if err != nil {
		t.Error(err)
	}
	if res1.Tenant.TenantId != res0.TenantId {
		t.Errorf("incorrect tenant id: %v, expected: %v", res1.Tenant.TenantId, res0.TenantId)
	}
	if res1.Tenant.Name != name0 {
		t.Errorf("incorrect name: %v, expected: %v", res1.Tenant.Name, name0)
	}

	_, err = esc.RenameTenant(t.Context(), &liara.RenameTenantRequest{
		TenantId: res0.TenantId,
		Name:     name1,
	})
	if err != nil {
		t.Error(err)
	}

	res3, err := esc.GetTenant(t.Context(), &liara.GetTenantRequest{
		TenantId: res0.TenantId,
	})
	if err != nil {
		t.Error(err)
	}
	if res3.Tenant.TenantId != res0.TenantId {
		t.Errorf("incorrect tenant id: %v, expected: %v", res3.Tenant.TenantId, res0.TenantId)
	}
	if res3.Tenant.Name != name1 {
		t.Errorf("incorrect name: %v, expected: %v", res3.Tenant.Name, name1)
	}

	_, err = esc.DeleteTenant(t.Context(), &liara.DeleteTenantRequest{
		TenantId: res0.TenantId,
	})
	if err != nil {
		t.Error(err)
	}

	_, err = esc.GetTenant(t.Context(), &liara.GetTenantRequest{
		TenantId: res0.TenantId,
	})
	if err == nil {
		t.Error("should return error")
	}
}
