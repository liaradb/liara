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
	kn.page.SetDirty()

	return i, true
}

func (kn *KeyNode) Insert(key Key, block BlockPosition) (iter.Seq[KeyEntry], iter.Seq[KeyEntry], bool) {
	ke := newKeyEntry(key, block)
	i := kn.searchIndex(ke.key)

	_, b, ok := kn.page.Insert(int16(ke.Size()), i)
	if !ok {
		// Split
		a, b := kn.split(i, ke)
		return a, b, false
	}

	ke.Write(b)
	kn.page.SetDirty()

	return nil, nil, true
}

func (kn *KeyNode) split(i int16, ke KeyEntry) (iter.Seq[KeyEntry], iter.Seq[KeyEntry]) {
	mid := kn.mid()
	return kn.first(i, mid, ke),
		kn.second(i, mid, ke)
}

func (kn *KeyNode) mid() int16 {
	return kn.page.Count() / 2
}

func (kn *KeyNode) first(i int16, mid int16, ke KeyEntry) func(yield func(KeyEntry) bool) {
	if i >= mid {
		return kn.childrenRange(0, mid)
	}

	return func(yield func(KeyEntry) bool) {
		if i == 0 {
			if !yield(ke) {
				return
			}
		}

		var j int16
		for e := range kn.childrenRange(0, mid) {
			if !yield(e) {
				return
			}

			j++

			if i == j {
				if !yield(ke) {
					return
				}
			}
		}
	}
}

func (kn *KeyNode) second(i int16, mid int16, ke KeyEntry) func(yield func(KeyEntry) bool) {
	if i < mid {
		return kn.childrenRange(mid, -1)
	}

	return func(yield func(KeyEntry) bool) {
		k := i - mid
		if k == 0 {
			if !yield(ke) {
				return
			}
		}

		var j int16
		for e := range kn.childrenRange(mid, -1) {
			if !yield(e) {
				return
			}

			j++

			if k == j {
				if !yield(ke) {
					return
				}
			}
		}
	}
}

// TODO: Test this
func (kn *KeyNode) Fill(entries iter.Seq[KeyEntry]) Key {
	var k Key
	first := true
	for e := range entries {
		if first {
			k = e.key
		}
		first = false
		// This will definitely fit
		_, _ = kn.Append(e.key, e.block)
	}
	kn.page.SetDirty()
	return k
}

// TODO: Test this
// TODO: Find a faster way
func (kn *KeyNode) Replace(entries iter.Seq[KeyEntry]) {
	cache := make([]KeyEntry, 0, kn.mid())
	for e := range entries {
		cache = append(cache, e)
	}

	kn.page.Clear()

	for _, e := range cache {
		// This will definitely fit
		_, _ = kn.Append(e.key, e.block)
	}

	kn.page.SetDirty()
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

func (kn *KeyNode) childrenRange(start, end int16) iter.Seq[KeyEntry] {
	return func(yield func(KeyEntry) bool) {
		for b := range kn.page.ChildrenRange(start, end) {
			le := KeyEntry{}
			le.Read(b)
			if !yield(le) {
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
