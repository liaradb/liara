package async

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
)

func TestCommand(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCommand)
}

func testCommand(t *testing.T) {
	h := make(CommandHandler[string])
	var errValue = errors.New("error value")

	go func() {
		if r, ok := h.Listen(t.Context()); ok {
			if v := r.Value(); v != "a" {
				t.Errorf("incorrect value: %v, expected: %v", v, "a")
			}

			r.Reply(errValue)
		}
	}()

	err := h.Send(t.Context(), "a")
	if !errors.Is(err, errValue) {
		t.Errorf("incorrect error: %v, expected: %v", err, errValue)
	}
}

func TestCommand_Listen__Canceled(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCommand_Listen__Canceled)
}

func testCommand_Listen__Canceled(t *testing.T) {
	h := make(CommandHandler[string])

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if _, ok := h.Listen(ctx); ok {
		t.Error("should return false")
	}
}

func TestCommand_Send__Canceled(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCommand_Send__Canceled)
}

func testCommand_Send__Canceled(t *testing.T) {
	h := make(CommandHandler[string])

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if err := h.Send(ctx, ""); err == nil {
		t.Error("should be canceled")
	}
}

func TestCommand__Wait__Canceled(t *testing.T) {
	t.Parallel()
	synctest.Test(t, testCommand__Wait__Canceled)
}

func testCommand__Wait__Canceled(t *testing.T) {
	h := make(CommandHandler[string])

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

	err := h.Send(ctx, "")
	if err == nil {
		t.Error("should be canceled")
	}
}
