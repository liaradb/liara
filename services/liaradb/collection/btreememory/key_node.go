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
	nodes    []node[K]
	leftID   storage.Offset
	rightID  storage.Offset
}

type keyEntry[K cmp.Ordered] struct {
	k  K
	id storage.Offset
}

var _ node[int] = (*keyNode[int])(nil)

func newKeyNode[K cmp.Ordered](s Storage[K], a, b node[K]) *keyNode[K] {
	kn := &keyNode[K]{
		i:       nextID(),
		storage: s,
		level:   a.height() + 1,
		children: []*keyEntry[K]{
			{k: a.key(), id: a.id()},
			{k: b.key(), id: b.id()}},
		nodes: []node[K]{a, b},
	}
	return kn
}

func (kn *keyNode[K]) key() K             { return kn.k }
func (kn *keyNode[K]) id() storage.Offset { return kn.i }

func (kn *keyNode[K]) count() int {
	count := 0
	for _, l := range kn.nodes {
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
	a := kn.nodes[0]

	l := len(kn.nodes)
	for i := 1; i < l; i++ {
		b := kn.nodes[i]
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
	kn.nodes = slices.Insert(kn.nodes, i, n)
	if len(kn.nodes) <= f {
		return nil, false
	}

	return kn.split(), true
}

func (kn *keyNode[K]) getInsertionIndex(k K) int {
	for i := len(kn.nodes) - 1; i >= 0; i-- {
		j := kn.nodes[i]
		if k >= j.key() {
			return i + 1
		}
	}
	return 0
}

func (kn *keyNode[K]) split() node[K] {
	half := len(kn.nodes) / 2

	kn2 := &keyNode[K]{
		i:        nextID(),
		k:        kn.nodes[half].key(),
		children: kn.children[half:],
		nodes:    kn.nodes[half:],
		leftID:   kn.i,
		rightID:  kn.rightID,
	}

	// TODO: Should we copy slices?
	kn.children = slices.Clone(kn.children[:half])
	kn.nodes = slices.Clone(kn.nodes[:half])
	kn.rightID = kn2.i

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
