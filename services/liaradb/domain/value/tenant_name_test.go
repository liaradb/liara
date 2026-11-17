package value

import (
	"testing"

	"github.com/liaradb/liaradb/encoder/buffer"
)

func TestTenantName(t *testing.T) {
	n := NewTenantName("name")

	b := buffer.New(n.Size())
	if err := n.Write(b); err != nil {
		t.Fatal(err)
	}

	b.SeekStart()

	r := TenantName{}
	if err := r.Read(b); err != nil {
		t.Fatal(err)
	}

	if r.String() != n.String() {
		t.Errorf("incorrect result: %v, expected: %v", r.String(), n.String())
	}
}
