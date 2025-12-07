package btree

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/encoder/raw"
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

	// TODO: Change to bool instead of error
	if err := ke.Write(raw.NewBufferFromSlice(b)); err != nil {
		return 0, false
	}

	return i, true
}

func (kn *KeyNode) Insert(key Key, block BlockPosition) (int16, bool) {
	ke := newKeyEntry(key, block)
	i, _ := kn.searchIndex(ke.key)

	i, b, ok := kn.page.Insert(int16(ke.Size()), i)
	if !ok {
		return 0, false
	}

	// TODO: Change to bool instead of error
	if err := ke.Write(raw.NewBufferFromSlice(b)); err != nil {
		return 0, false
	}

	return i, true
}

// TODO: Change to bool instead of error
func (kn *KeyNode) Children() iter.Seq2[KeyEntry, error] {
	return func(yield func(KeyEntry, error) bool) {
		for b := range kn.page.Children() {
			ke := KeyEntry{}
			if err := ke.Read(raw.NewBufferFromSlice(b)); err != nil {
				yield(KeyEntry{}, err)
				return
			}

			if !yield(ke, nil) {
				return
			}
		}
	}
}

// TODO: Change to bool instead of error
func (kn *KeyNode) Search(k Key) (BlockPosition, error) {
	p := BlockPosition(kn.page.LowID())
	for ke, err := range kn.Children() {
		if err != nil {
			return 0, err
		}

		if k < ke.key {
			break
		}

		p = ke.block
	}
	return p, nil
}

func (kn *KeyNode) searchIndex(k Key) (int16, error) {
	var i int16 = 0
	for ke, err := range kn.Children() {
		if err != nil {
			return 0, err
		}

		if k <= ke.key {
			break
		}

		i++
	}
	return i, nil
}
