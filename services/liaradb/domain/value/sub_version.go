package value

import "github.com/liaradb/liaradb/encoder/base"

const SubVersionSize = base.Uint32Size

type SubVersion struct {
	baseUint32
}

func NewSubVersion(value uint32) SubVersion {
	return SubVersion{baseUint32(value)}
}

func (sv SubVersion) Increment() SubVersion {
	return NewSubVersion(sv.Value() + 1)
}
