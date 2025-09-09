package log

import "encoding/binary"

type LogMagic uint32

var (
	LogMagicPage = LogMagic(binary.BigEndian.Uint32([]byte("PAGE")))
)

func (lm LogMagic) String() string {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(lm))
	return string(data)
}
