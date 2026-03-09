package value

import "github.com/liaradb/liaradb/encoder/base"

type Version struct {
	baseUint64
}

func NewVersion(value uint64) Version {
	return Version{baseUint64(value)}
}

func (v *Version) Increment() {
	v.baseUint64++
}

const VersionSize = base.BaseUint64Size
