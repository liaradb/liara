package async

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
)

func TestRequest(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testRequest)
}

func testRequest(t *testing.T) {
	h := make(Handler[string, int])
	var errValue = errors.New("error value")

	go func() {
		if r, ok := h.Listen(t.Context()); ok {
			if v := r.Value(); v != "a" {
				t.Errorf("incorrect value: %v, expected: %v", v, "a")
			}

			if c := r.Context(); c != t.Context() {
				t.Errorf("incorrect context: %v, expected: %v", c, t.Context())
			}

			r.Reply(2, errValue)
		}
	}()

	v, err := h.Send(t.Context(), "a")
	if v != 2 {
		t.Errorf("incorrect result: %v, expected: %v", v, 2)
	}
	if !errors.Is(err, errValue) {
		t.Errorf("incorrect error: %v, expected: %v", err, errValue)
	}
}

func TestRequest_Listen__Canceled(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testRequest_Listen__Canceled)
}

func testRequest_Listen__Canceled(t *testing.T) {
	h := make(Handler[string, string])

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if _, ok := h.Listen(ctx); ok {
		t.Error("should return false")
	}
}

func TestRequest_Send__Canceled(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testRequest_Send__Canceled)
}

func testRequest_Send__Canceled(t *testing.T) {
	h := make(Handler[string, string])

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if _, err := h.Send(ctx, ""); err == nil {
		t.Error("should be canceled")
	}
}

func TestRequest__Wait__Canceled(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testRequest__Wait__Canceled)
}

func testRequest__Wait__Canceled(t *testing.T) {
	h := make(Handler[string, string])

	ctx, cancel := context.WithCancel(t.Context())

	go func() {
		_, ok := h.Listen(ctx)
		if !ok {
			t.Error("should return true")
		}
	}()

	go func() {
		synctest.Wait()
		cancel()
	}()

	_, err := h.Send(ctx, "")
	if err == nil {
		t.Error("should be canceled")
	}
}
