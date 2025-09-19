package log

import "testing"

func TestSegmentName(t *testing.T) {
	t.Parallel()

	for message, test := range map[string]struct {
		id   SegmentID
		lsn  LogSequenceNumber
		name string
	}{
		"should handle index 0":         {0, 0, "segment_0000000000000000_0000000000000000.lr"},
		"should add index padding":      {1, 0, "segment_0000000000000001_0000000000000000.lr"},
		"should handle full size index": {234, 0, "segment_00000000000000ea_0000000000000000.lr"},
		"should handle lsn 0":           {0, 0, "segment_0000000000000000_0000000000000000.lr"},
		"should add lsn padding":        {0, 1, "segment_0000000000000000_0000000000000001.lr"},
		"should handle full size lsn":   {0, 234, "segment_0000000000000000_00000000000000ea.lr"},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()

			sn := NewSegmentName(test.id, test.lsn)
			if sn != ParseSegmentName(test.name) {
				t.Errorf("%v: incorrect value: %v, expected: %v", message, sn, test.name)
			}

			if i := sn.ID(); i != test.id {
				t.Errorf("%v: incorrect id: %v, expected: %v", message, i, test.id)
			}

			if l := sn.LogSequenceNumber(); l != test.lsn {
				t.Errorf("%v: incorrect log sequence number: %v, expected: %v", message, l, test.lsn)
			}
		})
	}
}

func TestSegmentName_Next(t *testing.T) {
	sn := NewSegmentName(1, 10)
	n := sn.Next(15)
	if i := n.ID(); i != 2 {
		t.Errorf("incorrect id: %v, expected: %v", i, 2)
	}
	if l := n.LogSequenceNumber(); l != 15 {
		t.Errorf("incorrect log sequence number: %v, expected: %v", l, 15)
	}
}
