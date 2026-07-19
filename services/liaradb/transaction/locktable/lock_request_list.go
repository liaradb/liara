package locktable

import (
	"container/list"
)

type lockRequestList[K comparable] struct {
	count  int
	shared int
	list   *list.List
}

func newLockRequestList[K comparable]() *lockRequestList[K] {
	return &lockRequestList[K]{
		list: list.New(),
	}
}

func (lrl *lockRequestList[K]) push(r *lockRequest[K]) {
	lrl.list.PushFront(r)
}

func (lrl *lockRequestList[K]) pop() (*lockRequest[K], bool) {
	lrl.decrement()

	if lrl.isShared() {
		lrl.unShare()
	} else if lrl.isExcluded() {
		lrl.unExclude()
	}

	e := lrl.list.Front()
	if e == nil {
		return nil, false
	}

	return lrl.list.Remove(e).(*lockRequest[K]), true
}

func (lrl *lockRequestList[K]) lock(lr *lockRequest[K]) {
	switch lr.lockRequestType() {
	case lockRequestTypeS:
		lrl.sLock(lr)
	case lockRequestTypeX:
		lrl.xLock(lr)
	case lockRequestTypeUpgrade:
		lrl.upgradeLock(lr)
	}
}

func (lrl *lockRequestList[K]) xLock(lr *lockRequest[K]) {
	lrl.increment()

	if lrl.isFirst() {
		lrl.exclude()
		lr.send()
	} else {
		lrl.push(lr)
	}
}

func (lrl *lockRequestList[K]) upgradeLock(lr *lockRequest[K]) {
	if lrl.isFirst() || lrl.isSharedOnce() {
		lrl.exclude()
		lr.send()
	} else {
		lrl.push(lr)
	}
}

func (lrl *lockRequestList[K]) sLock(lr *lockRequest[K]) {
	lrl.increment()

	if lrl.isFirst() || lrl.isShared() {
		lrl.share()
		lr.send()
	} else {
		lrl.push(lr)
	}
}

func (lrl *lockRequestList[K]) unlock() (result bool) {
	for {
		lr, ok := lrl.pop()
		if !ok || lrl.sendUnlock(lr) {
			return lrl.isEmpty()
		}
	}
}

func (lrl *lockRequestList[K]) sendUnlock(lr *lockRequest[K]) bool {
	if !lr.send() {
		return false
	}

	if lr.isShared() {
		lrl.share()
	} else {
		lrl.exclude()
	}

	return true
}

func (lrl *lockRequestList[K]) isShared() bool     { return lrl.shared > 0 }
func (lrl *lockRequestList[K]) isSharedOnce() bool { return lrl.shared == 1 }
func (lrl *lockRequestList[K]) isExcluded() bool   { return lrl.shared < 0 }
func (lrl *lockRequestList[K]) isEmpty() bool      { return lrl.count == 0 }
func (lrl *lockRequestList[K]) isFirst() bool      { return lrl.count == 1 }

func (lrl *lockRequestList[K]) share()     { lrl.shared++ }
func (lrl *lockRequestList[K]) unShare()   { lrl.shared-- }
func (lrl *lockRequestList[K]) exclude()   { lrl.shared = -1 }
func (lrl *lockRequestList[K]) unExclude() { lrl.shared = 0 }
func (lrl *lockRequestList[K]) increment() { lrl.count++ }
func (lrl *lockRequestList[K]) decrement() { lrl.count-- }
