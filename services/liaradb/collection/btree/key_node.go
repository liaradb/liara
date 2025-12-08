package btree

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/page"
)

type KeyNode struct {
	page page.BTreePage
}

func newKeyNode(page page.BTreePage) *KeyNode {
	return &KeyNode{
		page: page,
	}
}

func (kn *KeyNode) Init(p BlockPosition) {
	kn.page.SetLowID(p.Value())
}

// TODO: Test this
func (kn *KeyNode) Append(key Key, block BlockPosition) (int16, bool) {
	ke := newKeyEntry(key, block)
	i, b, ok := kn.page.Append(int16(ke.Size()))
	if !ok {
		return 0, false
	}

	ke.Write(b)

	return i, true
}

func (kn *KeyNode) Insert(key Key, block BlockPosition) (int16, bool) {
	ke := newKeyEntry(key, block)
	i := kn.searchIndex(ke.key)

	i, b, ok := kn.page.Insert(int16(ke.Size()), i)
	if !ok {
		return 0, false
	}

	ke.Write(b)

	return i, true
}

func (kn *KeyNode) Children() iter.Seq[KeyEntry] {
	return func(yield func(KeyEntry) bool) {
		for b := range kn.page.Children() {
			ke := KeyEntry{}
			ke.Read(b)
			if !yield(ke) {
				return
			}
		}
	}
}

// TODO: Handle not found
func (kn *KeyNode) Search(k Key) BlockPosition {
	p := BlockPosition(kn.page.LowID())
	for ke := range kn.Children() {
		if k < ke.key {
			break
		}

		p = ke.block
	}
	return p
}

// TODO: Handle not found
func (kn *KeyNode) searchIndex(k Key) int16 {
	var i int16 = 0
	for ke := range kn.Children() {
		if k <= ke.key {
			break
		}

		i++
	}
	return i
}
