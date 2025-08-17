package locktable

import (
	"context"
	"fmt"
)

type lockType int8

const (
	lockTypeNone lockType = iota
	lockTypeS
	lockTypeX
)

type ConcurrencyMgr[K comparable] struct {
	lockTable *LockTable[K]
	locks     map[K]lockType
}

func NewConcurrencyMgr[K comparable](lt *LockTable[K]) *ConcurrencyMgr[K] {
	return &ConcurrencyMgr[K]{
		lockTable: lt,
		locks:     make(map[K]lockType),
	}
}

func (cm *ConcurrencyMgr[K]) SLock(ctx context.Context, blk K) error {
	if cm.hasAnyLock(blk) {
		return nil
	}

	if ok := cm.lockTable.sLock(ctx, blk); !ok {
		return fmt.Errorf("unable to SLock %v: %w", blk, ErrLockAbort)
	}

	cm.locks[blk] = lockTypeS

	return nil
}

func (cm *ConcurrencyMgr[K]) XLock(ctx context.Context, blk K) error {
	switch cm.locks[blk] {
	case lockTypeX:
		return nil
	case lockTypeS:
		if ok := cm.lockTable.upgradeLock(ctx, blk); !ok {
			return fmt.Errorf("unable to UpgradeLock %v: %w", blk, ErrLockAbort)
		}
	default:
		if ok := cm.lockTable.xLock(ctx, blk); !ok {
			return fmt.Errorf("unable to XLock %v: %w", blk, ErrLockAbort)
		}
	}

	cm.locks[blk] = lockTypeX

	return nil
}

func (cm *ConcurrencyMgr[K]) Release() {
	for blk := range cm.locks {
		cm.lockTable.release(blk)
	}

	clear(cm.locks)
}

func (cm *ConcurrencyMgr[K]) hasAnyLock(blk K) bool {
	l := cm.locks[blk]
	return l > lockTypeNone
}
