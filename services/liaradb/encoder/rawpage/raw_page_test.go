package rawpage

import "testing"

func TestRawPage(t *testing.T) {
	p := New(make([]byte, 256))
	if _, ok := p.Append(16); !ok {
		t.Error("should get a buffer")
	}
}
