package value

import "github.com/liaradb/liaradb/encoder/base"

const GlobalVersionSize = base.Uint64Size

type GlobalVersion struct {
	baseUint64
}

func NewGlobalVersion(value uint64) GlobalVersion {
	return GlobalVersion{baseUint64(value)}
}
