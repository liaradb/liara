package fixedv2

import "github.com/liaradb/liaradb/storage"

type bufferSlice []*storage.Buffer

func (bs *bufferSlice) Append(b *storage.Buffer) {
	*bs = append(*bs, b)
}

func (bs *bufferSlice) Release() {
	for _, b := range *bs {
		b.Release()
	}
}
