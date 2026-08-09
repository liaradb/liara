package bufferpage

import (
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/storage"
)

const (
	FragmentHeaderSize = base.Uint16Size +
		base.Uint16Size +
		page.CrcSize
)

type BufferPage struct {
	*page.Page
	header
	buffer *storage.Buffer
}

func (bp *BufferPage) Clear() {
	bp.Page.Clear()
	bp.header.init()
}

func New(b *storage.Buffer) *BufferPage {
	page := page.NewFromSlice(b.Raw(), headerSize, FragmentHeaderSize)
	header, _ := newHeader(page.Header())
	return &BufferPage{
		Page:   page,
		header: header,
		buffer: b,
	}
}

func (bp *BufferPage) Fill(data []byte) {
	bp.Page.Fill(data)
}

func (bp *BufferPage) Shadow(base *BufferPage) {
	bp.Page.Fill(base.Data())
}

func (bp *BufferPage) Release() {
	bp.buffer.Release()
}
