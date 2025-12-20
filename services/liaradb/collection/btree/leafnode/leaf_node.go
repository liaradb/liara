package leafnode

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/collection/btree/value"
)

type LeafNode struct {
	page page.Page
}

type Iterator = iter.Seq2[value.Key, value.RecordID]

func New(page page.Page) *LeafNode {
	return &LeafNode{
		page: page,
	}
}

func (ln *LeafNode) LeftID() value.BlockPosition {
	return value.BlockPosition(ln.page.LowID())
}

func (ln *LeafNode) RightID() value.BlockPosition {
	return value.BlockPosition(ln.page.HighID())
}

// TODO: Test this
func (ln *LeafNode) SetLeftID(block value.BlockPosition) {
	ln.page.SetLowID(block.Value())
	ln.page.SetDirty()
}

// TODO: Test this
func (ln *LeafNode) SetRightID(block value.BlockPosition) {
	ln.page.SetHighID(block.Value())
	ln.page.SetDirty()
}

func (ln *LeafNode) Append(key value.Key, recordID value.RecordID) (int16, bool) {
	le := newLeafEntry(key, recordID)
	i, b, ok := ln.page.Append(int16(le.Size()))
	if !ok {
		return 0, false
	}

	le.Write(b)
	ln.page.SetDirty()

	return i, true
}

func (ln *LeafNode) Insert(key value.Key, recordID value.RecordID) (Iterator, Iterator, bool) {
	le := newLeafEntry(key, recordID)
	i := ln.searchIndexRange(le.key)

	_, b, ok := ln.page.Insert(int16(le.Size()), i)
	if !ok {
		a, b := ln.split(i, le)
		return a, b, false
	}

	le.Write(b)
	ln.page.SetDirty()

	return nil, nil, true
}

// TODO: Test this
func (ln *LeafNode) Fill(
	leftID value.BlockPosition,
	rightID value.BlockPosition,
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
	ln.page.SetDirty()
	return k
}

// TODO: Test this
// TODO: Find a faster way
func (ln *LeafNode) Replace(rightID value.BlockPosition, entries Iterator) {
	cache := make([]leafEntry, 0, ln.mid())
	for key, rid := range entries {
		cache = append(cache, newLeafEntry(key, rid))
	}

	leftID := ln.LeftID()

	ln.page.Clear()

	for _, e := range cache {
		// This will definitely fit
		_, _ = ln.Append(e.key, e.recordID)
	}

	// TODO: We are duplicating set dirty calls
	ln.SetLeftID(leftID)
	ln.SetRightID(rightID)
	ln.page.SetDirty()
}

func (ln *LeafNode) split(i int16, le leafEntry) (Iterator, Iterator) {
	mid := ln.mid()
	return ln.first(i, mid, le), ln.second(i, mid, le)
}

func (ln *LeafNode) mid() int16 {
	return ln.page.Count() / 2
}

func (ln *LeafNode) first(i int16, mid int16, le leafEntry) Iterator {
	if i >= mid {
		return ln.childrenRange(0, mid)
	}

	// TODO: Simplify this
	return func(yield func(value.Key, value.RecordID) bool) {
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
	return func(yield func(value.Key, value.RecordID) bool) {
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
	b, ok := ln.page.Child(index)
	if !ok {
		return leafEntry{}, false
	}

	le := leafEntry{}
	le.Read(b)

	return le, true
}

func (ln *LeafNode) Children() Iterator {
	return func(yield func(value.Key, value.RecordID) bool) {
		for b := range ln.page.Children() {
			le := leafEntry{}
			le.Read(b)
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}
	}
}

func (ln *LeafNode) childrenRange(start, end int16) Iterator {
	return func(yield func(value.Key, value.RecordID) bool) {
		for b := range ln.page.ChildrenRange(start, end) {
			le := leafEntry{}
			le.Read(b)
			if !yield(le.Key(), le.RecordID()) {
				return
			}
		}
	}
}

func (ln *LeafNode) RecordIDs() iter.Seq[value.RecordID] {
	return func(yield func(value.RecordID) bool) {
		for _, rid := range ln.Children() {
			if !yield(rid) {
				return
			}
		}
	}
}

func (ln *LeafNode) Search(k value.Key) (value.RecordID, bool) {
	i, ok := ln.searchIndex(k)
	if !ok {
		return value.RecordID{}, false
	}

	le, ok := ln.Child(i)
	if !ok {
		return value.RecordID{}, false
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
func (ln *LeafNode) Release()  { ln.page.Release() }
func (ln *LeafNode) Latch()    { ln.page.Latch() }
func (ln *LeafNode) Unlatch()  { ln.page.Unlatch() }
func (ln *LeafNode) RLatch()   { ln.page.RLatch() }
func (ln *LeafNode) RUnlatch() { ln.page.RUnlatch() }

func (ln *LeafNode) SearchRange(k value.Key) iter.Seq[value.RecordID] {
	return func(yield func(value.RecordID) bool) {
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
