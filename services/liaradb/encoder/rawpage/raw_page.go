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

func (p *RawPage) Append(size int32) (int32, *raw.Buffer, bool) {
	offset := p.list.Next() - size
	i, ok := p.list.Push(offset)
	if !ok {
		return 0, nil, false
	}

	p.list.SetNext(offset)

	b, ok := p.byteList.Slice(int64(offset), int64(size))
	if !ok {
		return 0, nil, false
	}

	return i, b, true
}

func (p RawPage) Length() int32 {
	return int32(len(p.data))
}

func (p RawPage) Space() int32 {
	return p.list.Next() - p.list.Size()
}
