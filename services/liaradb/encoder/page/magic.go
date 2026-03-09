package page

import (
	"encoding/binary"
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type Magic uint32

const MagicSize = 4

var (
	MagicEmpty = Magic(0)
	MagicPage  = Magic(binary.BigEndian.Uint32([]byte("PAGE")))
	MagicFree  = Magic(binary.BigEndian.Uint32([]byte("FREE")))
)

func (m Magic) Write(w io.Writer) error {
	return raw.WriteInt32(w, m)
}

func (m *Magic) Read(r io.Reader) error {
	if err := m.read(r); err != nil {
		return err
	}

	return m.Validate()
}

func (m *Magic) read(r io.Reader) error {
	return raw.ReadInt32(r, m)
}

func (m Magic) Validate() error {
	switch m {
	case MagicEmpty:
		return nil
	case MagicFree:
		return nil
	case MagicPage:
		return nil
	default:
		return ErrNotPage
	}
}

func (m Magic) IsEmpty() bool {
	return m == MagicEmpty
}

func (m Magic) IsPage() bool {
	return m == MagicPage
}

func (m Magic) String() string {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(m))
	return string(data)
}
