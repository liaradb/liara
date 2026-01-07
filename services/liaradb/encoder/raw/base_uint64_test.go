package raw

import "testing"

func TestBaseUint64_String(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value BaseUint64
		want  string
	}{
		"should handle 0": {
			value: NewBaseUint64(0),
			want:  "0000000000000000",
		},
		"should handle 1": {
			value: NewBaseUint64(1),
			want:  "0000000000000001",
		},
		"should handle 2": {
			value: NewBaseUint64(2),
			want:  "0000000000000002",
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			if s := c.value.String(); s != c.want {
				t.Errorf("%v: incorrect string: %v, expected: %v", message, s, c.want)
			}
		})
	}
}
