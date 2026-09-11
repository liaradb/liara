package leafnode

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/storage/link"
)

type LeafNode struct {
	node node.Node
}

type Iterator = iter.Seq2[key.Key, link.RecordLocator]

func New(page node.Node) *LeafNode {
	return &LeafNode{
		node: page,
	}
}

func (ln *LeafNode) LeftID() link.FilePosition {
	return ln.node.LowID()
}

func (ln *LeafNode) RightID() link.FilePosition {
	return ln.node.HighID()
}

func (ln *LeafNode) SetLeftID(block link.FilePosition) {
	ln.setLeftID(block)
	ln.node.SetDirty()
}

func (ln *LeafNode) setLeftID(block link.FilePosition) {
	ln.node.SetLowID(block)
}

func (ln *LeafNode) SetRightID(block link.FilePosition) {
	ln.setRightID(block)
	ln.node.SetDirty()
}

func (ln *LeafNode) setRightID(block link.FilePosition) {
	ln.node.SetHighID(block)
}

func (ln *LeafNode) Append(key key.Key, recordID link.RecordLocator) bool {
	le := newLeafEntry(key, recordID)
	b, ok := ln.node.Append(int16(le.Size()))
	if !ok {
		return false
	}

	le.Write(b)
	ln.node.SetDirty()

	return true
}

func (ln *LeafNode) Insert(key key.Key, recordID link.RecordLocator) (Iterator, Iterator, bool) {
	le := newLeafEntry(key, recordID)
	i := ln.searchIndexRange(le.key)

	b, ok := ln.node.Insert(int16(le.Size()), i)
	if !ok {
		a, b := ln.split(i, le)
		return a, b, false
	}

	le.Write(b)
	ln.node.SetDirty()

	return nil, nil, true
}

func (ln *LeafNode) Fill(
	leftID link.FilePosition,
	rightID link.FilePosition,
	entries Iterator,
) key.Key {
	var k key.Key
	first := true
	for key, rid := range entries {
		if first {
			k = key
		}
		first = false
		// This will definitely fit
		_ = ln.Append(key, rid)
	}

	ln.setLeftID(leftID)
	ln.setRightID(rightID)
	ln.node.SetDirty()
	return k
}

// TODO: Find a faster way
func (ln *LeafNode) Replace(rightID link.FilePosition, entries Iterator) {
	cache := make([]leafEntry, 0, ln.mid())
	for key, rid := range entries {
		cache = append(cache, newLeafEntry(key, rid))
	}

	leftID := ln.LeftID()

	ln.node.Clear()

	for _, e := range cache {
		// This will definitely fit
		_ = ln.Append(e.key, e.recordID)
	}

	ln.setLeftID(leftID)
	ln.setRightID(rightID)
	ln.node.SetDirty()
}

func (ln *LeafNode) split(i link.SlotID, le leafEntry) (Iterator, Iterator) {
	mid := ln.mid()
	return ln.first(i, mid, le), ln.second(i, mid, le)
}

func (ln *LeafNode) mid() link.SlotID {
	return ln.node.Count() / 2
}

func (ln *LeafNode) first(i, mid link.SlotID, le leafEntry) Iterator {
	// Iterate the first half
	if i >= mid {
		return ln.childrenRange(0, mid)
	}

	return func(yield func(key.Key, link.RecordLocator) bool) {
		// i is 0
		if i == 0 {
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}

		var j link.SlotID
		for key, rid := range ln.childrenRange(0, mid) {
			if !yield(key, rid) {
				return
			}

			j++

			// i is greater than 0
			if i == j {
				if !yield(le.Key(), le.RecordID()) {
					return
				}
			}
		}
	}
}

func (ln *LeafNode) second(i, mid link.SlotID, le leafEntry) Iterator {
	// Iterate the second half
	if i < mid {
		return ln.childrenRange(mid, -1)
	}

	return func(yield func(key.Key, link.RecordLocator) bool) {
		k := i - mid
		// k is 0
		if k == 0 {
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}

		var j link.SlotID
		for key, rid := range ln.childrenRange(mid, -1) {
			if !yield(key, rid) {
				return
			}

			j++

			// k is greater than 0
			if k == j {
				if !yield(le.Key(), le.RecordID()) {
					return
				}
			}
		}
	}
}

func (ln *LeafNode) Child(index link.SlotID) (leafEntry, bool) {
	b, ok := ln.node.Child(index)
	if !ok {
		return leafEntry{}, false
	}

	le := leafEntry{}
	le.Read(b)

	return le, true
}

func (ln *LeafNode) Children() Iterator {
	return func(yield func(key.Key, link.RecordLocator) bool) {
		for b := range ln.node.Children() {
			le := leafEntry{}
			le.Read(b)
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}
	}
}

func (ln *LeafNode) childrenRange(start, end link.SlotID) Iterator {
	return func(yield func(key.Key, link.RecordLocator) bool) {
		for b := range ln.node.ChildrenRange(start, end) {
			le := leafEntry{}
			le.Read(b)
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}
	}
}

func (ln *LeafNode) RecordIDs() iter.Seq[link.RecordLocator] {
	return func(yield func(link.RecordLocator) bool) {
		for _, rid := range ln.Children() {
			if !yield(rid) {
				return
			}
		}
	}
}

func (ln *LeafNode) Search(k key.Key) (link.RecordLocator, bool) {
	i, ok := ln.searchIndex(k)
	if !ok {
		return link.RecordLocator{}, false
	}

	le, ok := ln.Child(i)
	if !ok {
		return link.RecordLocator{}, false
	}

	return le.recordID, true
}

func (ln *LeafNode) searchIndex(k key.Key) (link.SlotID, bool) {
	var i link.SlotID = 0
	for key := range ln.Children() {
		if k.Equal(key) {
			return i, true
		}
		if k.LessEqual(key) {
			return 0, false
		}

		i++
	}
	return 0, false
}

func (ln *LeafNode) searchIndexRange(k key.Key) link.SlotID {
	var i link.SlotID = 0
	for key := range ln.Children() {
		if k.LessEqual(key) {
			break
		}

		i++
	}
	return i
}

// TODO: Test this
func (ln *LeafNode) Release()  { ln.node.Release() }
func (ln *LeafNode) Latch()    { ln.node.Latch() }
func (ln *LeafNode) Unlatch()  { ln.node.Unlatch() }
func (ln *LeafNode) RLatch()   { ln.node.RLatch() }
func (ln *LeafNode) RUnlatch() { ln.node.RUnlatch() }

func (ln *LeafNode) SearchRange(k key.Key) iter.Seq[link.RecordLocator] {
	return func(yield func(link.RecordLocator) bool) {
		i := ln.searchIndexRange(k)
		for _, rid := range ln.childrenRange(i, -1) {
			if !yield(rid) {
				return
			}
		}
	}
}
