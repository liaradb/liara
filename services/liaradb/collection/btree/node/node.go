package node

import (
	"context"
	"iter"
	"slices"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/encoder/page"
	"github.com/liaradb/liaradb/recovery/logpage"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

const (
	itemSize = 4
)

type Node struct {
	header
	page    *page.Page
	buffer  *storage.Buffer
	actions []action
}

type action interface {
	Log(context.Context, span.Log) (logpage.LogSequenceNumber, error)
}

func New(buffer *storage.Buffer) Node {
	page := page.NewFromSlice(buffer.Raw(), HeaderSize, 0)
	header, _ := newHeader(page.Header())

	return Node{
		header: header,
		page:   page,
		buffer: buffer,
	}
}

func (n *Node) Clear() {
	// n.buffer.Clear()
	n.addAction(newClearAction(n.buffer.FileBlockID().RecordID(0), n.page.Data()))
	n.page.Clear()
}

// TODO: Test this
func (n *Node) Release()  { n.buffer.Release() }
func (n *Node) Latch()    { n.buffer.Latch() }
func (n *Node) Unlatch()  { n.buffer.Unlatch() }
func (n *Node) RLatch()   { n.buffer.RLatch() }
func (n *Node) RUnlatch() { n.buffer.RUnlatch() }

func (n *Node) SetDirty() {
	n.buffer.SetDirty()
}

func (n *Node) SetLevel(l byte) {
	n.header.setLevel(l)
}

func (n *Node) LogSequenceNumber() logpage.LogSequenceNumber {
	return logpage.NewLogSequenceNumber(n.page.TrackingID())
}

func (n *Node) setLogSequenceNumber(lsn logpage.LogSequenceNumber) {
	n.page.SetTrackingID(lsn.Value())
}

func (n *Node) addAction(a action) {
	n.actions = append(n.actions, a)
}

func (n *Node) LogActions(ctx context.Context, l span.Log) error {
	for _, a := range n.actions {
		lsn, err := a.Log(ctx, l)
		if err != nil {
			return err
		}

		n.setLogSequenceNumber(lsn)
	}

	return nil
}

func (n *Node) Append(size int16) ([]byte, bool) {
	// TODO: Fix this cast
	_, b, ok := n.page.NextMustFit(int(size))
	return b, ok
}

func (n *Node) CommitAppend(size int16) {
	// TODO: Fix this cast
	_, b, _ := n.page.NextMustFit(int(size))
	sid := n.page.CommitAppend(int(size))
	rid := n.buffer.FileBlockID().RecordID(sid)

	n.addAction(newAppendAction(rid, b))
	n.addHeaderAction(rid)
}

func (n *Node) CommitInsert(size int16, index link.SlotID) {
	// TODO: Fix this cast
	_, b, _ := n.page.NextMustFit(int(size))
	sid := n.page.Insert(int(size), index)
	rid := n.buffer.FileBlockID().RecordID(sid)

	n.addAction(newInsertAction(rid, b))
	n.addHeaderAction(rid)
}

func (n *Node) addHeaderAction(rid link.RecordID) {
	if n.header.dirty {
		n.addAction(newHeaderAction(rid,
			slices.Clone(n.header.data),
			slices.Clone(n.header.base)))
		n.header.commit()
	}
}

func (n *Node) Count() link.SlotID {
	return n.page.Count()
}

func (n *Node) Space() int16 {
	return int16(n.page.Space())
}

func (n *Node) Child(index link.SlotID) []byte {
	_, b := n.page.Slot(index)
	return b
}

func (n *Node) Children() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for _, b := range n.page.Slots() {
			if !yield(b) {
				return
			}
		}
	}
}

func (n *Node) ChildrenRange(start, end link.SlotID) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for _, b := range n.page.SlotsRange(start, end) {
			if !yield(b) {
				return
			}
		}
	}
}
