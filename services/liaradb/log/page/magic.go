package page

import (
	"encoding/binary"
	"io"

	"github.com/liaradb/liaradb/raw"
)

type Magic uint32

const magicSize = 4

var (
	MagicPage = Magic(binary.BigEndian.Uint32([]byte("PAGE")))
)

func (m Magic) Write(w io.Writer) error {
	return raw.WriteInt32(w, m)
}

func (m *Magic) Read(r io.Reader) error {
	return raw.ReadInt32(r, m)
}

func (m *Magic) ReadIsPage(r io.Reader) error {
	var b Magic
	if err := b.Read(r); err != nil {
		return err
	}

	if b != MagicPage {
		return ErrNotPage
	}

	return nil
}

func (m Magic) String() string {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(m))
	return string(data)
}
