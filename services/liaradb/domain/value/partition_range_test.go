package value

import "testing"

func TestParitionRange(t *testing.T) {
	for message, c := range map[string]struct {
		skip bool
		a    PartitionID
		b    PartitionID
		low  PartitionID
		high PartitionID
	}{
		"should create": {
			a:    NewPartitionID(1),
			b:    NewPartitionID(2),
			low:  NewPartitionID(1),
			high: NewPartitionID(2),
		},
		"should create in reverse": {
			a:    NewPartitionID(2),
			b:    NewPartitionID(1),
			low:  NewPartitionID(1),
			high: NewPartitionID(2),
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			pr := NewPartitionRange(c.a, c.b)

			if l := pr.Low(); l != c.low {
				t.Errorf("incorrect value: %v, expected: %v", l, c.low)
			}

			if h := pr.High(); h != c.high {
				t.Errorf("incorrect value: %v, expected: %v", h, c.high)
			}

			if l, h := pr.All(); l != c.low {
				t.Errorf("incorrect value: %v, expected: %v", l, c.low)
			} else if h != c.high {
				t.Errorf("incorrect value: %v, expected: %v", h, c.high)
			}
		})
	}
}

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
