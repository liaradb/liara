package value

import (
	"encoding/binary"

	"github.com/liaradb/liaradb/encoder/base"
)

const GlobalVersionSize = base.Uint64Size

type GlobalVersion struct {
	baseUint64
}

func NewGlobalVersion(value uint64) GlobalVersion {
	return GlobalVersion{baseUint64(value)}
}

func (gv GlobalVersion) Increment() GlobalVersion {
	return NewGlobalVersion(gv.Value() + 1)
}

func (gv GlobalVersion) Bytes() []byte {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], gv.Value())
	return data[:]
}
