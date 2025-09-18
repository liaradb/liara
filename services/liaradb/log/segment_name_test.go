package log

import "testing"

func TestSegmentName(t *testing.T) {
	t.Parallel()

	for message, test := range map[string]struct {
		index SegmentID
		lsn   LogSequenceNumber
		name  string
	}{
		"should handle index 0":         {0, 0, "segment_000_000.lr"},
		"should add index padding":      {1, 0, "segment_001_000.lr"},
		"should handle full size index": {234, 0, "segment_234_000.lr"},
		"should overflow index":         {1234, 0, "segment_1234_000.lr"},
		"should handle lsn 0":           {0, 0, "segment_000_000.lr"},
		"should add lsn padding":        {0, 1, "segment_000_001.lr"},
		"should handle full size lsn":   {0, 234, "segment_000_234.lr"},
		"should overflow osn":           {0, 1234, "segment_000_1234.lr"},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()

			sn := NewSegmentName(test.index, test.lsn)
			if sn != ParseSegmentName(test.name) {
				t.Errorf("%v: incorrect value: %v, expected: %v", message, sn, test.name)
			}

			if i := sn.ID(); i != test.index {
				t.Errorf("%v: incorrect index: %v, expected: %v", message, i, test.index)
			}

			if l := sn.LogSequenceNumber(); l != test.lsn {
				t.Errorf("%v: incorrect log sequence number: %v, expected: %v", message, l, test.lsn)
			}
		})
	}
}
