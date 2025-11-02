package btree

import (
	"cmp"
	"io"

	"github.com/liaradb/liaradb/encoder/page"
)

type nodeItem[K cmp.Ordered, V any] struct {
	node[K, V]
}

var _ page.ItemSerializer = (*nodeItem[int, int])(nil)

func (n *nodeItem[K, V]) Read(io.Reader, page.CRC) error {
	return nil
}

func (n *nodeItem[K, V]) Size() int {
	return 0
}

func (n *nodeItem[K, V]) Write(io.Writer) (page.CRC, error) {
	return page.CRC{}, nil
}
