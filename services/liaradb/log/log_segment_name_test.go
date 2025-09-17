package log

import "testing"

func TestLogSegmentName(t *testing.T) {
	t.Parallel()

	for message, test := range map[string]struct {
		index int
		name  string
	}{
		"should handle 0":         {0, "segment_000.lr"},
		"should add padding":      {1, "segment_001.lr"},
		"should handle full size": {234, "segment_234.lr"},
		"should overflow":         {1234, "segment_1234.lr"},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()

			lsn := NewLogSegmentName(test.index)
			if lsn != LogSegmentName(test.name) {
				t.Errorf("%v: incorrect value: %v, expected: %v", message, lsn, test.name)
			}
		})
	}
}
