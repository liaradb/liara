package btree

import "github.com/liaradb/liaradb/encoder/page"

type KeyEntry[K Key] struct {
	Key      K
	Position page.Offset
}
