package keynode

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/storage/link"
)

type KeyNode struct {
	node node.Node
}

type Iterator = iter.Seq2[key.Key, link.FilePosition]

func New(page node.Node) *KeyNode {
	return &KeyNode{
		node: page,
	}
}

func (kn *KeyNode) append(key key.Key, block link.FilePosition) (int16, bool) {
	ke := newKeyEntry(key, block)
	i, b, ok := kn.node.Append(int16(ke.Size()))
	if !ok {
		return 0, false
	}

	ke.Write(b)
	kn.node.SetDirty()

	return i, true
}

func (kn *KeyNode) Insert(key key.Key, block link.FilePosition) (Iterator, Iterator, bool) {
	ke := newKeyEntry(key, block)
	i := kn.searchIndex(ke.key)

	_, b, ok := kn.node.Insert(int16(ke.Size()), i)
	if !ok {
		// Split
		a, b := kn.split(i, ke)
		return a, b, false
	}

	ke.Write(b)
	kn.node.SetDirty()

	return nil, nil, true
}

func (kn *KeyNode) split(i int16, ke keyEntry) (Iterator, Iterator) {
	mid := kn.mid()
	return kn.first(i, mid, ke),
		kn.second(i, mid, ke)
}

func (kn *KeyNode) mid() int16 {
	return kn.node.Count() / 2
}

func (kn *KeyNode) first(i int16, mid int16, ke keyEntry) Iterator {
	if i >= mid {
		return kn.childrenRange(0, mid)
	}

	return func(yield func(key.Key, link.FilePosition) bool) {
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

	return func(yield func(key.Key, link.FilePosition) bool) {
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
func (kn *KeyNode) Fill(l byte, entries Iterator) key.Key {
	var k key.Key
	first := true
	for key, block := range entries {
		if first {
			k = key
		}
		first = false
		// This will definitely fit
		_, _ = kn.append(key, block)
	}

	kn.node.SetLevel(l)
	kn.node.SetDirty()
	return k
}

// TODO: Test this
// TODO: Find a faster way
func (kn *KeyNode) Replace(l byte, entries Iterator) {
	cache := make([]keyEntry, 0, kn.mid())
	for key, block := range entries {
		cache = append(cache, newKeyEntry(key, block))
	}

	kn.node.Clear()

	for _, e := range cache {
		// This will definitely fit
		_, _ = kn.append(e.key, e.block)
	}

	kn.node.SetLevel(l)
	kn.node.SetDirty()
}

// TODO: Test this
func (kn *KeyNode) ReplaceRoot(l byte, block0 link.FilePosition, key1 key.Key, block1 link.FilePosition) bool {
	// This should always have a child
	// TODO: Will this always be the lower key?
	key0, _, _ := kn.Child(0)

	kn.node.Clear()

	if _, ok := kn.append(key0, block0); !ok {
		return false
	}

	_, ok := kn.append(key1, block1)
	kn.node.SetLevel(l)
	kn.node.SetDirty()
	return ok
}

func (kn *KeyNode) Children() Iterator {
	return func(yield func(key.Key, link.FilePosition) bool) {
		for b := range kn.node.Children() {
			ke := keyEntry{}
			ke.Read(b)
			if !yield(ke.Key(), ke.Block()) {
				return
			}
		}
	}
}

// TODO: Test this
func (kn *KeyNode) Child(i int16) (key.Key, link.FilePosition, bool) {
	b, ok := kn.node.Child(i)
	if !ok {
		return key.Key{}, 0, false
	}

	ke := keyEntry{}
	ke.Read(b)
	return ke.Key(), ke.Block(), true
}

func (kn *KeyNode) childrenRange(start, end int16) Iterator {
	return func(yield func(key.Key, link.FilePosition) bool) {
		for b := range kn.node.ChildrenRange(start, end) {
			ke := keyEntry{}
			ke.Read(b)
			if !yield(ke.Key(), ke.Block()) {
				return
			}
		}
	}
}

func (kn *KeyNode) Search(k key.Key) link.FilePosition {
	var p link.FilePosition
	first := true
	for key, block := range kn.Children() {
		if first {
			p = block
			first = false
			continue
		}
		if k.Less(key) {
			break
		}

		p = block
	}
	return p
}

func (kn *KeyNode) searchIndex(k key.Key) int16 {
	var i int16 = 0
	for key := range kn.Children() {
		if k.LessEqual(key) {
			break
		}

		i++
	}
	return i
}

func (kn *KeyNode) Level() byte {
	return kn.node.Level()
}

// TODO: Test this
func (kn *KeyNode) Release()  { kn.node.Release() }
func (kn *KeyNode) Latch()    { kn.node.Latch() }
func (kn *KeyNode) Unlatch()  { kn.node.Unlatch() }
func (kn *KeyNode) RLatch()   { kn.node.RLatch() }
func (kn *KeyNode) RUnlatch() { kn.node.RUnlatch() }
