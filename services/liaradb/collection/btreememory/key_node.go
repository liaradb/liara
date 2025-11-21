package btreememory

import (
	"cmp"
	"slices"

	"github.com/liaradb/liaradb/storage"
)

type keyNode[K cmp.Ordered] struct {
	storage  Storage[K]
	i        storage.Offset
	k        K
	level    int
	children []*keyEntry[K]
	leftID   storage.Offset
	rightID  storage.Offset
}

type keyEntry[K cmp.Ordered] struct {
	k  K
	id storage.Offset
}

var _ node[int] = (*keyNode[int])(nil)

func newKeyNode[K cmp.Ordered](s Storage[K], i storage.Offset, a, b node[K]) *keyNode[K] {
	kn := &keyNode[K]{
		i:       i,
		storage: s,
		level:   a.height() + 1,
		children: []*keyEntry[K]{
			{k: a.key(), id: a.id()},
			{k: b.key(), id: b.id()}},
	}
	return kn
}

func (kn *keyNode[K]) key() K             { return kn.k }
func (kn *keyNode[K]) id() storage.Offset { return kn.i }
func (kn *keyNode[K]) isKeyNode() bool    { return true }
func (kn *keyNode[K]) isLeafNode() bool   { return false }
func (kn *keyNode[K]) count() int         { return len(kn.children) }
func (kn *keyNode[K]) height() int        { return kn.level }

func (kn *keyNode[K]) getChild(k K) (storage.Offset, bool) {
	a := kn.children[0]

	l := len(kn.children)
	for i := 1; i < l; i++ {
		b := kn.children[i]
		if k < b.k {
			return a.id, true
		}

		a = b
	}

	return a.id, true
}

func (kn *keyNode[K]) getValue(k K) (RecordID, bool) {
	return RecordID{}, false
}

func (kn *keyNode[K]) insert(f int, k K, id storage.Offset) (*keyNode[K], bool) {
	i := kn.getInsertionIndex(k)
	if i == 0 {
		kn.k = k
	}

	// TODO: Split before inserting
	kn.children = slices.Insert(kn.children, i, &keyEntry[K]{k: k, id: id})
	if len(kn.children) <= f {
		return nil, false
	}

	return kn.split(), true
}

func (kn *keyNode[K]) getInsertionIndex(k K) int {
	for i := len(kn.children) - 1; i >= 0; i-- {
		j := kn.children[i]
		if k >= j.k {
			return i + 1
		}
	}
	return 0
}

func (kn *keyNode[K]) split() *keyNode[K] {
	half := len(kn.children) / 2

	kn2 := &keyNode[K]{
		i:        kn.storage.NextID(),
		k:        kn.children[half].k,
		children: kn.children[half:],
		leftID:   kn.i,
		rightID:  kn.rightID,
	}

	// TODO: Should we copy slices?
	kn.children = slices.Clone(kn.children[:half])
	kn.rightID = kn2.i

	return kn2
}

func (kn *keyNode[K]) delete(f int, k K, rid RecordID) {

}

func (kn *keyNode[K]) deleteAll(f int, k K) {

}
