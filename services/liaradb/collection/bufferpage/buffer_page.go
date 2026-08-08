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

func (lp *BufferPage) Clear() {
	lp.Page.Clear()
	lp.header.init()
}

func New(b *storage.Buffer) *BufferPage {
	page := page.NewFromSlice(b.Raw(), headerSize, FragmentHeaderSize)
	header, _ := newHeader(page.Header())
	return &BufferPage{
		Page:   page,
		header: header,
	}
}

func (lp *BufferPage) Fill(data []byte) {
	lp.Page.Fill(data)
}

func (lp *BufferPage) Shadow(base *BufferPage) {
	lp.Page.Fill(base.Data())
}
