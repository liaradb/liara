package page

import (
	"io"
	"testing"

	"github.com/liaradb/liaradb/encoder/buffer"
	"github.com/liaradb/liaradb/util/testing/iotesting"
)

func TestMagic(t *testing.T) {
	t.Parallel()

	r, w := iotesting.NewReaderWriter()

	var m Magic = MagicPage
	if err := m.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var m2 Magic
	if err := m2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if m != m2 {
		t.Errorf("incorrect value: %v, expected: %v", m2, m)
	}
}

func TestMagic_Validate(t *testing.T) {
	for message, c := range map[string]struct {
		skip  bool
		m     Magic
		valid bool
	}{
		"should handle empty": {
			valid: true,
		},
		"should handle free": {
			m:     MagicFree,
			valid: true,
		},
		"should handle page": {
			m:     MagicPage,
			valid: true,
		},
		"should return error on unknown": {
			m:     Magic(10),
			valid: false,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			err := c.m.Validate()
			if c.valid {
				if err != nil {
					t.Error(err)
				}
			} else {
				if err != ErrNotPage {
					t.Error("should return error")
				}
			}
		})
	}
}

func TestMagic_Read(t *testing.T) {
	for message, c := range map[string]struct {
		skip  bool
		b     *buffer.Buffer
		valid bool
	}{
		"should handle empty": {
			b:     buffer.New(4),
			valid: true,
		},
		"should handle free": {
			b:     buffer.NewFromSlice([]byte(MagicFree.String())),
			valid: true,
		},
		"should handle page": {
			b:     buffer.NewFromSlice([]byte(MagicPage.String())),
			valid: true,
		},
		"should return error on unknown": {
			b:     buffer.NewFromSlice([]byte("abcd")),
			valid: false,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			var m Magic
			err := m.Read(c.b)
			if c.valid {
				if err != nil {
					t.Error(err)
				}
			} else {
				if err != ErrNotPage {
					t.Error("should return error")
				}
			}
		})
	}
}

func TestMagic_Read__Error(t *testing.T) {
	b := buffer.New(0)
	var m Magic
	if err := m.Read(b); err == nil {
		t.Error("should return error")
	}
}

func TestMagicEmpty(t *testing.T) {
	t.Parallel()

	if s := MagicEmpty.String(); s != "\x00\x00\x00\x00" {
		t.Error("value is incorrect")
	}

	if !MagicEmpty.IsEmpty() {
		t.Error("should be empty")
	}

	if MagicEmpty.IsPage() {
		t.Error("should not be page")
	}
}

func TestMagicFree(t *testing.T) {
	t.Parallel()

	if s := MagicFree.String(); s != "FREE" {
		t.Error("value is incorrect")
	}

	if MagicFree.IsEmpty() {
		t.Error("should not be empty")
	}

	if MagicFree.IsPage() {
		t.Error("should not be page")
	}
}

func TestMagicPage(t *testing.T) {
	t.Parallel()

	if s := MagicPage.String(); s != "PAGE" {
		t.Error("value is incorrect")
	}

	if MagicPage.IsEmpty() {
		t.Error("should not be empty")
	}

	if !MagicPage.IsPage() {
		t.Error("should be page")
	}
}
