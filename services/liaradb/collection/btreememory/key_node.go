package btreememory

import (
	"cmp"
	"slices"

	"github.com/liaradb/liaradb/storage"
)

type keyNode[K cmp.Ordered] struct {
	storage  Storage[K]
	id       storage.Offset
	k        K
	level    int
	children []node[K]
	leftID   storage.Offset
	rightID  storage.Offset
}

var _ node[int] = (*keyNode[int])(nil)

func newKeyNode[K cmp.Ordered](s Storage[K], a, b node[K]) *keyNode[K] {
	kn := &keyNode[K]{
		id:       nextID(),
		storage:  s,
		level:    a.height() + 1,
		children: []node[K]{a, b},
	}
	return kn
}

func (kn *keyNode[K]) key() K {
	return kn.k
}

func (kn *keyNode[K]) count() int {
	count := 0
	for _, l := range kn.children {
		count += l.count()
	}
	return count
}

func (kn *keyNode[K]) getValue(k K) (RecordID, bool) {
	if kn == nil || kn.count() == 0 {
		return kn.zero()
	}

	return kn.getChild(k).getValue(k)
}

func (kn *keyNode[K]) getChild(k K) node[K] {
	a := kn.children[0]

	l := len(kn.children)
	for i := 1; i < l; i++ {
		b := kn.children[i]
		if k < b.key() {
			return a
		}

		a = b
	}

	return a
}

func (kn *keyNode[K]) insert(f int, k K, rid RecordID) (node[K], bool) {
	n, ok := kn.getChild(k).insert(f, k, rid)
	if !ok {
		return nil, false
	}

	return kn.insertNode(f, k, n)
}

func (kn *keyNode[K]) insertNode(f int, k K, n node[K]) (node[K], bool) {
	i := kn.getInsertionIndex(n.key())
	if i == 0 {
		kn.k = k
	}

	// TODO: Split before inserting
	kn.children = slices.Insert(kn.children, i, n)
	if len(kn.children) <= f {
		return nil, false
	}

	return kn.split(), true
}

func (kn *keyNode[K]) getInsertionIndex(k K) int {
	for i := len(kn.children) - 1; i >= 0; i-- {
		j := kn.children[i]
		if k >= j.key() {
			return i + 1
		}
	}
	return 0
}

func (kn *keyNode[K]) split() node[K] {
	half := len(kn.children) / 2

	kn2 := &keyNode[K]{
		id:       nextID(),
		k:        kn.children[half].key(),
		children: kn.children[half:],
		leftID:   kn.id,
		rightID:  kn.rightID,
	}

	// TODO: Should we copy slices?
	kn.children = slices.Clone(kn.children[:half])
	kn.rightID = kn2.id

	return kn2
}

func (kn *keyNode[K]) delete(f int, k K, rid RecordID) {

}

func (kn *keyNode[K]) deleteAll(f int, k K) {

}

func (kn *keyNode[K]) height() int {
	if kn == nil {
		return 0
	}

	return kn.level
}

func (*keyNode[K]) zero() (RecordID, bool) {
	return RecordID{}, false
}
