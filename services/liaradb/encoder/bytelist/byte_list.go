package bytelist

import (
	"github.com/liaradb/liaradb/encoder/raw"
)

type ByteList struct {
	data []byte
}

func New(data []byte) ByteList {
	return ByteList{
		data: data,
	}
}

func (l ByteList) Slice(off int64, n int64) (*raw.Buffer, bool) {
	if off >= int64(len(l.data)) {
		return nil, false
	}

	end := off + n
	if end > int64(len(l.data)) {
		return nil, false
	}

	return raw.NewBufferFromSlice(l.data[off:end]), true
}
