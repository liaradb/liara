package btree

import (
	"cmp"

	"github.com/liaradb/liaradb/encoder/page"
)

type Key interface {
	cmp.Ordered
	page.Serializer
}
