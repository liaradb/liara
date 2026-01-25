package value

import "github.com/liaradb/liaradb/encoder/raw"

type Version struct {
	baseUint64
}

func NewVersion(value uint64) Version {
	return Version{baseUint64(value)}
}

const VersionSize = raw.BaseUint64Size
