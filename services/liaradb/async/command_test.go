package async

import (
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
		r := <-h
		r.Reply(errValue)
	}()

	err := h.Send(t.Context(), "a")
	if !errors.Is(err, errValue) {
		t.Errorf("incorrect error: %v, expected: %v", err, errValue)
	}
}
