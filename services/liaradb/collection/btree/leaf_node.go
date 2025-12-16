package btree

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/key"
	"github.com/liaradb/liaradb/collection/btree/page"
)

type LeafNode struct {
	page page.BTreePage
}

func newLeafNode(page page.BTreePage) *LeafNode {
	return &LeafNode{
		page: page,
	}
}

func (ln *LeafNode) LeftID() BlockPosition {
	return BlockPosition(ln.page.LowID())
}

func (ln *LeafNode) RightID() BlockPosition {
	return BlockPosition(ln.page.HighID())
}

func (ln *LeafNode) Append(key key.Key, recordID RecordID) (int16, bool) {
	le := newLeafEntry(key, recordID)
	i, b, ok := ln.page.Append(int16(le.Size()))
	if !ok {
		return 0, false
	}

	le.Write(b)
	ln.page.SetDirty()

	return i, true
}

func (ln *LeafNode) Insert(key key.Key, recordID RecordID) (iter.Seq[LeafEntry], iter.Seq[LeafEntry], bool) {
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
func (ln *LeafNode) Fill(entries iter.Seq[LeafEntry]) key.Key {
	var k key.Key
	first := true
	for e := range entries {
		if first {
			k = e.key
		}
		first = false
		// This will definitely fit
		_, _ = ln.Append(e.key, e.recordID)
	}
	ln.page.SetDirty()
	return k
}

// TODO: Test this
// TODO: Find a faster way
func (ln *LeafNode) Replace(entries iter.Seq[LeafEntry]) {
	cache := make([]LeafEntry, 0, ln.mid())
	for e := range entries {
		cache = append(cache, e)
	}

	ln.page.Clear()

	for _, e := range cache {
		// This will definitely fit
		_, _ = ln.Append(e.key, e.recordID)
	}

	ln.page.SetDirty()
}

func (ln *LeafNode) split(i int16, le LeafEntry) (iter.Seq[LeafEntry], iter.Seq[LeafEntry]) {
	mid := ln.mid()
	return ln.first(i, mid, le), ln.second(i, mid, le)
}

func (ln *LeafNode) mid() int16 {
	return ln.page.Count() / 2
}

func (ln *LeafNode) first(i int16, mid int16, le LeafEntry) func(yield func(LeafEntry) bool) {
	if i >= mid {
		return ln.childrenRange(0, mid)
	}

	return func(yield func(LeafEntry) bool) {
		if i == 0 {
			if !yield(le) {
				return
			}
		}

		var j int16
		for e := range ln.childrenRange(0, mid) {
			if !yield(e) {
				return
			}

			j++

			if i == j {
				if !yield(le) {
					return
				}
			}
		}
	}
}

func (ln *LeafNode) second(i int16, mid int16, le LeafEntry) func(yield func(LeafEntry) bool) {
	if i < mid {
		return ln.childrenRange(mid, -1)
	}

	return func(yield func(LeafEntry) bool) {
		k := i - mid
		if k == 0 {
			if !yield(le) {
				return
			}
		}

		var j int16
		for e := range ln.childrenRange(mid, -1) {
			if !yield(e) {
				return
			}

			j++

			if k == j {
				if !yield(le) {
					return
				}
			}
		}
	}
}

func (ln *LeafNode) Child(index int16) (LeafEntry, bool) {
	b, ok := ln.page.Child(index)
	if !ok {
		return LeafEntry{}, false
	}

	le := LeafEntry{}
	le.Read(b)

	return le, true
}

func (ln *LeafNode) Children() iter.Seq[LeafEntry] {
	return func(yield func(LeafEntry) bool) {
		for b := range ln.page.Children() {
			le := LeafEntry{}
			le.Read(b)
			if !yield(le) {
				return
			}
		}
	}
}

func (ln *LeafNode) childrenRange(start, end int16) iter.Seq[LeafEntry] {
	return func(yield func(LeafEntry) bool) {
		for b := range ln.page.ChildrenRange(start, end) {
			le := LeafEntry{}
			le.Read(b)
			if !yield(le) {
				return
			}
		}
	}
}

func (ln *LeafNode) Search(k key.Key) (RecordID, bool) {
	i, ok := ln.searchIndex(k)
	if !ok {
		return RecordID{}, false
	}

	le, ok := ln.Child(i)
	if !ok {
		return RecordID{}, false
	}

	return le.recordID, true
}

func (ln *LeafNode) searchIndex(k key.Key) (int16, bool) {
	var i int16 = 0
	for ke := range ln.Children() {
		if k == ke.key {
			return i, true
		}
		if k <= ke.key {
			return 0, false
		}

		i++
	}
	return 0, false
}

func (ln *LeafNode) searchIndexRange(k key.Key) int16 {
	var i int16 = 0
	for ke := range ln.Children() {
		if k <= ke.key {
			break
		}

		i++
	}
	return i
}
