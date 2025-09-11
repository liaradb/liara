package log

import (
	"encoding/binary"
	"io"
)

type LogMagic uint32

var (
	LogMagicPage = LogMagic(binary.BigEndian.Uint32([]byte("PAGE")))
)

func (lm *LogMagic) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, *lm)
}

func (lm *LogMagic) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, lm)
}

func (lm *LogMagic) String() string {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(*lm))
	return string(data)
}
