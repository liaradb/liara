package value

import "testing"

func TestPartitionRange_ReadWrite(t *testing.T) {
	t.Parallel()

	pr := NewPartitionRange(NewPartitionID(2), NewPartitionID(3))

	data := make([]byte, 12)
	data0 := pr.WriteData(data)

	if l := len(data0); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	pr2 := NewPartitionRange(NewPartitionID(0), NewPartitionID(0))
	data1 := pr2.ReadData(data)

	if l := len(data1); l != 4 {
		t.Errorf("incorrect length: %v, expected: %v", l, 4)
	}

	if v := pr2.Low().Value(); v != 2 {
		t.Errorf("incorrect value: %v, expected: %v", v, 2)
	}

	if v := pr2.High().Value(); v != 3 {
		t.Errorf("incorrect value: %v, expected: %v", v, 3)
	}
}
