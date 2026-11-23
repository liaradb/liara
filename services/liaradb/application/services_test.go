package application

import "testing"

func TestCollections(t *testing.T) {
	c := newCollections()
	if c == nil {
		t.Error("should have value")
	}
	if c.EventLog == nil {
		t.Error("should have value")
	}
	if c.Idempotency == nil {
		t.Error("should have value")
	}
	if c.Outbox == nil {
		t.Error("should have value")
	}
}
