package value

import "testing"

func TestTenantName(t *testing.T) {
	n := NewTenantName("name")
	data := make([]byte, TenantNameSize)
	_ = n.WriteData(data)

	r := TenantName{}
	r.ReadData(data)
	if r.String() != n.String() {
		t.Errorf("incorrect result: %v, expected: %v", r.String(), n.String())
	}
}
