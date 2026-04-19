package storage

import "testing"

func TestBufferStatus_String(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip   bool
		bs     BufferStatus
		result string
	}{
		"should handle unknown": {
			bs:     100,
			result: "unknown",
		},
		"should handle uninitialized": {
			bs:     BufferStatusUninitialized,
			result: "uninitialized",
		},
		"should handle corrupt": {
			bs:     BufferStatusCorrupt,
			result: "corrupt",
		},
		"should handle dirty": {
			bs:     BufferStatusDirty,
			result: "dirty",
		},
		"should handle loaded": {
			bs:     BufferStatusLoaded,
			result: "loaded",
		},
		"should handle loading": {
			bs:     BufferStatusLoading,
			result: "loading",
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			if s := c.bs.String(); s != c.result {
				t.Errorf("incorrect string: %v, expected: %v", s, c.result)
			}
		})
	}
}
