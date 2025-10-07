package async

import (
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
		r := <-h
		r.Reply(2, errValue)
	}()

	v, err := h.Send(t.Context(), "a")
	if v != 2 {
		t.Errorf("incorrect result: %v, expected: %v", v, 2)
	}
	if !errors.Is(err, errValue) {
		t.Errorf("incorrect error: %v, expected: %v", err, errValue)
	}
}
