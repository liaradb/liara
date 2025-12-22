package leafnode

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/node"
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/encoder/page"
)

type LeafNode struct {
	node node.Node
}

type Iterator = iter.Seq2[value.Key, value.RecordLocator]

func New(page node.Node) *LeafNode {
	return &LeafNode{
		node: page,
	}
}

func (ln *LeafNode) LeftID() page.Offset {
	return page.Offset(ln.node.LowID())
}

func (ln *LeafNode) RightID() page.Offset {
	return page.Offset(ln.node.HighID())
}

// TODO: Test this
func (ln *LeafNode) SetLeftID(block page.Offset) {
	ln.node.SetLowID(block.Value())
	ln.node.SetDirty()
}

// TODO: Test this
func (ln *LeafNode) SetRightID(block page.Offset) {
	ln.node.SetHighID(block.Value())
	ln.node.SetDirty()
}

func (ln *LeafNode) Append(key value.Key, recordID value.RecordLocator) (int16, bool) {
	le := newLeafEntry(key, recordID)
	i, b, ok := ln.node.Append(int16(le.Size()))
	if !ok {
		return 0, false
	}

	le.Write(b)
	ln.node.SetDirty()

	return i, true
}

func (ln *LeafNode) Insert(key value.Key, recordID value.RecordLocator) (Iterator, Iterator, bool) {
	le := newLeafEntry(key, recordID)
	i := ln.searchIndexRange(le.key)

	_, b, ok := ln.node.Insert(int16(le.Size()), i)
	if !ok {
		a, b := ln.split(i, le)
		return a, b, false
	}

	le.Write(b)
	ln.node.SetDirty()

	return nil, nil, true
}

// TODO: Test this
func (ln *LeafNode) Fill(
	leftID page.Offset,
	rightID page.Offset,
	entries Iterator,
) value.Key {
	var k value.Key
	first := true
	for key, rid := range entries {
		if first {
			k = key
		}
		first = false
		// This will definitely fit
		_, _ = ln.Append(key, rid)
	}

	// TODO: We are duplicating set dirty calls
	ln.SetLeftID(leftID)
	ln.SetRightID(rightID)
	ln.node.SetDirty()
	return k
}

// TODO: Test this
// TODO: Find a faster way
func (ln *LeafNode) Replace(rightID page.Offset, entries Iterator) {
	cache := make([]leafEntry, 0, ln.mid())
	for key, rid := range entries {
		cache = append(cache, newLeafEntry(key, rid))
	}

	leftID := ln.LeftID()

	ln.node.Clear()

	for _, e := range cache {
		// This will definitely fit
		_, _ = ln.Append(e.key, e.recordID)
	}

	// TODO: We are duplicating set dirty calls
	ln.SetLeftID(leftID)
	ln.SetRightID(rightID)
	ln.node.SetDirty()
}

func (ln *LeafNode) split(i int16, le leafEntry) (Iterator, Iterator) {
	mid := ln.mid()
	return ln.first(i, mid, le), ln.second(i, mid, le)
}

func (ln *LeafNode) mid() int16 {
	return ln.node.Count() / 2
}

func (ln *LeafNode) first(i int16, mid int16, le leafEntry) Iterator {
	if i >= mid {
		return ln.childrenRange(0, mid)
	}

	// TODO: Simplify this
	return func(yield func(value.Key, value.RecordLocator) bool) {
		if i == 0 {
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}

		var j int16
		for key, rid := range ln.childrenRange(0, mid) {
			if !yield(key, rid) {
				return
			}

			j++

			if i == j {
				if !yield(le.Key(), le.RecordID()) {
					return
				}
			}
		}
	}
}

func (ln *LeafNode) second(i int16, mid int16, le leafEntry) Iterator {
	if i < mid {
		return ln.childrenRange(mid, -1)
	}

	// TODO: Simplify this
	return func(yield func(value.Key, value.RecordLocator) bool) {
		k := i - mid
		if k == 0 {
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}

		var j int16
		for key, rid := range ln.childrenRange(mid, -1) {
			if !yield(key, rid) {
				return
			}

			j++

			if k == j {
				if !yield(le.Key(), le.RecordID()) {
					return
				}
			}
		}
	}
}

func (ln *LeafNode) Child(index int16) (leafEntry, bool) {
	b, ok := ln.node.Child(index)
	if !ok {
		return leafEntry{}, false
	}

	le := leafEntry{}
	le.Read(b)

	return le, true
}

func (ln *LeafNode) Children() Iterator {
	return func(yield func(value.Key, value.RecordLocator) bool) {
		for b := range ln.node.Children() {
			le := leafEntry{}
			le.Read(b)
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}
	}
}

func (ln *LeafNode) childrenRange(start, end int16) Iterator {
	return func(yield func(value.Key, value.RecordLocator) bool) {
		for b := range ln.node.ChildrenRange(start, end) {
			le := leafEntry{}
			le.Read(b)
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}
	}
}

func (ln *LeafNode) RecordIDs() iter.Seq[value.RecordLocator] {
	return func(yield func(value.RecordLocator) bool) {
		for _, rid := range ln.Children() {
			if !yield(rid) {
				return
			}
		}
	}
}

func (ln *LeafNode) Search(k value.Key) (value.RecordLocator, bool) {
	i, ok := ln.searchIndex(k)
	if !ok {
		return value.RecordLocator{}, false
	}

	le, ok := ln.Child(i)
	if !ok {
		return value.RecordLocator{}, false
	}

	return le.recordID, true
}

func (ln *LeafNode) searchIndex(k value.Key) (int16, bool) {
	var i int16 = 0
	for key := range ln.Children() {
		if k == key {
			return i, true
		}
		if k <= key {
			return 0, false
		}

		i++
	}
	return 0, false
}

func (ln *LeafNode) searchIndexRange(k value.Key) int16 {
	var i int16 = 0
	for key := range ln.Children() {
		if k <= key {
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

func (ln *LeafNode) SearchRange(k value.Key) iter.Seq[value.RecordLocator] {
	return func(yield func(value.RecordLocator) bool) {
		i, ok := ln.searchIndex(k)
		if !ok {
			return
		}

		for _, rid := range ln.childrenRange(i, -1) {
			if !yield(rid) {
				return
			}
		}
	}
}
