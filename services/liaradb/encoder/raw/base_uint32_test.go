package raw

import "testing"

func TestBaseUint32_String(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		skip  bool
		value BaseUint32
		want  string
	}{
		"should handle 0": {
			value: NewBaseUint32(0),
			want:  "00000000",
		},
		"should handle 1": {
			value: NewBaseUint32(1),
			want:  "00000001",
		},
		"should handle 2": {
			value: NewBaseUint32(2),
			want:  "00000002",
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
