package keynode

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/page"
)

type KeyNode struct {
	page page.BTreePage
}

func New(page page.BTreePage) *KeyNode {
	return &KeyNode{
		page: page,
	}
}

// TODO: Test this
func (kn *KeyNode) Append(key key.Key, block BlockPosition) (int16, bool) {
	ke := newKeyEntry(key, block)
	i, b, ok := kn.page.Append(int16(ke.Size()))
	if !ok {
		return 0, false
	}

	ke.Write(b)
	kn.page.SetDirty()

	return i, true
}

func (kn *KeyNode) Insert(key key.Key, block BlockPosition) (iter.Seq[KeyEntry], iter.Seq[KeyEntry], bool) {
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
func (kn *KeyNode) Fill(l byte, entries iter.Seq[KeyEntry]) key.Key {
	var k key.Key
	first := true
	for e := range entries {
		if first {
			k = e.key
		}
		first = false
		// This will definitely fit
		_, _ = kn.Append(e.key, e.block)
	}

	kn.page.SetLevel(l)
	kn.page.SetDirty()
	return k
}

// TODO: Test this
// TODO: Find a faster way
func (kn *KeyNode) Replace(l byte, entries iter.Seq[KeyEntry]) {
	cache := make([]KeyEntry, 0, kn.mid())
	for e := range entries {
		cache = append(cache, e)
	}

	kn.page.Clear()

	for _, e := range cache {
		// This will definitely fit
		_, _ = kn.Append(e.key, e.block)
	}

	kn.page.SetLevel(l)
	kn.page.SetDirty()
}

// TODO: Test this
func (kn *KeyNode) ReplaceRoot(l byte, block0 BlockPosition, key1 key.Key, block1 BlockPosition) bool {
	// This should always have a child
	// TODO: Will this always be the lower key?
	child0, _ := kn.Child(0)

	kn.page.Clear()

	if _, ok := kn.Append(child0.Key(), block0); !ok {
		return false
	}

	_, ok := kn.Append(key1, block1)
	kn.page.SetLevel(l)
	kn.page.SetDirty()
	return ok
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

// TODO: Test this
func (kn *KeyNode) Child(i int16) (KeyEntry, bool) {
	b, ok := kn.page.Child(i)
	if !ok {
		return KeyEntry{}, false
	}

	ke := KeyEntry{}
	ke.Read(b)
	return ke, true
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

func (kn *KeyNode) Search(k key.Key) BlockPosition {
	var p BlockPosition
	first := true
	for ke := range kn.Children() {
		if first {
			p = ke.block
			first = false
			continue
		}
		if k < ke.key {
			break
		}

		p = ke.block
	}
	return p
}

func (kn *KeyNode) searchIndex(k key.Key) int16 {
	var i int16 = 0
	for ke := range kn.Children() {
		if k <= ke.key {
			break
		}

		i++
	}
	return i
}

func (kn *KeyNode) Level() byte {
	return kn.page.Level()
}

// TODO: Test this
func (kn *KeyNode) Release() {
	kn.page.Release()
}
