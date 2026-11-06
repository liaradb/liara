package span

import "github.com/liaradb/liaradb/encoder/page"

const (
	headerSize = 0
)

type Page struct {
	*page.Page
}

func NewPage(data []byte) Page {
	return Page{page.NewFromSlice(data, headerSize, fragmentHeaderSize)}
}
