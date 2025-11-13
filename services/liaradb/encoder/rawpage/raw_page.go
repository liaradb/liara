package rawpage

import (
	"github.com/liaradb/liaradb/encoder/bytelist"
	"github.com/liaradb/liaradb/encoder/list"
	"github.com/liaradb/liaradb/encoder/raw"
)

type RawPage struct {
	data     []byte
	list     list.List
	byteList bytelist.ByteList
}

func New(data []byte) RawPage {
	return RawPage{
		data:     data,
		list:     list.New(data),
		byteList: bytelist.New(data),
	}
}

func (p RawPage) Append(size int32) (*raw.Buffer, bool) {
	offset := p.length() - size
	return p.byteList.Slice(int64(offset), int64(size))
}

func (p RawPage) length() int32 {
	return int32(len(p.data))
}
