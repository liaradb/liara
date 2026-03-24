package value

import "github.com/liaradb/liaradb/encoder/base"

type Version struct {
	baseUint64
}

func NewVersion(value int64) Version {
	return Version{baseUint64(value)}
}

func (v *Version) Increment() {
	v.baseUint64++
}

func (v Version) Value() int64 { return v.Signed() }

const VersionSize = base.Uint64Size
