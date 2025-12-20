package keynode

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/collection/btree/value"
)

type KeyNode struct {
	page page.Page
}

type Iterator = iter.Seq2[value.Key, BlockPosition]

func New(page page.Page) *KeyNode {
	return &KeyNode{
		page: page,
	}
}

// TODO: Test this
func (kn *KeyNode) Append(key value.Key, block BlockPosition) (int16, bool) {
	ke := newKeyEntry(key, block)
	i, b, ok := kn.page.Append(int16(ke.Size()))
	if !ok {
		return 0, false
	}

	ke.Write(b)
	kn.page.SetDirty()

	return i, true
}

func (kn *KeyNode) Insert(key value.Key, block BlockPosition) (Iterator, Iterator, bool) {
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

func (kn *KeyNode) split(i int16, ke keyEntry) (Iterator, Iterator) {
	mid := kn.mid()
	return kn.first(i, mid, ke),
		kn.second(i, mid, ke)
}

func (kn *KeyNode) mid() int16 {
	return kn.page.Count() / 2
}

func (kn *KeyNode) first(i int16, mid int16, ke keyEntry) Iterator {
	if i >= mid {
		return kn.childrenRange(0, mid)
	}

	return func(yield func(value.Key, BlockPosition) bool) {
		if i == 0 {
			if !yield(ke.Key(), ke.Block()) {
				return
			}
		}

		var j int16
		for key, block := range kn.childrenRange(0, mid) {
			if !yield(key, block) {
				return
			}

			j++

			if i == j {
				if !yield(ke.Key(), ke.Block()) {
					return
				}
			}
		}
	}
}

func (kn *KeyNode) second(i int16, mid int16, ke keyEntry) Iterator {
	if i < mid {
		return kn.childrenRange(mid, -1)
	}

	return func(yield func(value.Key, BlockPosition) bool) {
		k := i - mid
		if k == 0 {
			if !yield(ke.Key(), ke.Block()) {
				return
			}
		}

		var j int16
		for key, block := range kn.childrenRange(mid, -1) {
			if !yield(key, block) {
				return
			}

			j++

			if k == j {
				if !yield(ke.Key(), ke.Block()) {
					return
				}
			}
		}
	}
}

// TODO: Test this
func (kn *KeyNode) Fill(l byte, entries Iterator) value.Key {
	var k value.Key
	first := true
	for key, block := range entries {
		if first {
			k = key
		}
		first = false
		// This will definitely fit
		_, _ = kn.Append(key, block)
	}

	kn.page.SetLevel(l)
	kn.page.SetDirty()
	return k
}

// TODO: Test this
// TODO: Find a faster way
func (kn *KeyNode) Replace(l byte, entries Iterator) {
	cache := make([]keyEntry, 0, kn.mid())
	for key, block := range entries {
		cache = append(cache, newKeyEntry(key, block))
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
func (kn *KeyNode) ReplaceRoot(l byte, block0 BlockPosition, key1 value.Key, block1 BlockPosition) bool {
	// This should always have a child
	// TODO: Will this always be the lower key?
	key0, _, _ := kn.Child(0)

	kn.page.Clear()

	if _, ok := kn.Append(key0, block0); !ok {
		return false
	}

	_, ok := kn.Append(key1, block1)
	kn.page.SetLevel(l)
	kn.page.SetDirty()
	return ok
}

func (kn *KeyNode) Children() Iterator {
	return func(yield func(value.Key, BlockPosition) bool) {
		for b := range kn.page.Children() {
			ke := keyEntry{}
			ke.Read(b)
			if !yield(ke.Key(), ke.Block()) {
				return
			}
		}
	}
}

// TODO: Test this
func (kn *KeyNode) Child(i int16) (value.Key, BlockPosition, bool) {
	b, ok := kn.page.Child(i)
	if !ok {
		return "", 0, false
	}

	ke := keyEntry{}
	ke.Read(b)
	return ke.Key(), ke.Block(), true
}

func (kn *KeyNode) childrenRange(start, end int16) Iterator {
	return func(yield func(value.Key, BlockPosition) bool) {
		for b := range kn.page.ChildrenRange(start, end) {
			ke := keyEntry{}
			ke.Read(b)
			if !yield(ke.Key(), ke.Block()) {
				return
			}
		}
	}
}

func (kn *KeyNode) Search(k value.Key) BlockPosition {
	var p BlockPosition
	first := true
	for key, block := range kn.Children() {
		if first {
			p = block
			first = false
			continue
		}
		if k < key {
			break
		}

		p = block
	}
	return p
}

func (kn *KeyNode) searchIndex(k value.Key) int16 {
	var i int16 = 0
	for key := range kn.Children() {
		if k <= key {
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
func (kn *KeyNode) Release()  { kn.page.Release() }
func (kn *KeyNode) Latch()    { kn.page.Latch() }
func (kn *KeyNode) Unlatch()  { kn.page.Unlatch() }
func (kn *KeyNode) RLatch()   { kn.page.RLatch() }
func (kn *KeyNode) RUnlatch() { kn.page.RUnlatch() }
