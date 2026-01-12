package value

import "github.com/liaradb/liaradb/encoder/raw"

const GlobalVersionSize = raw.BaseUint64Size

type GlobalVersion struct {
	baseUint64
}

func NewGlobalVersion(value uint64) GlobalVersion {
	return GlobalVersion{baseUint64(value)}
}
