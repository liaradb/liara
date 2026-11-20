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

func (cm *ConcurrencyMgr[K]) SLock(ctx context.Context, k K) error {
	if cm.hasAnyLock(k) {
		return nil
	}

	if ok := cm.lockTable.sLock(ctx, k); !ok {
		return fmt.Errorf("unable to SLock %v: %w", k, ErrLockAbort)
	}

	cm.locks[k] = lockTypeS

	return nil
}

func (cm *ConcurrencyMgr[K]) XLock(ctx context.Context, k K) error {
	switch cm.locks[k] {
	case lockTypeX:
		return nil
	case lockTypeS:
		if ok := cm.lockTable.upgradeLock(ctx, k); !ok {
			return fmt.Errorf("unable to UpgradeLock %v: %w", k, ErrLockAbort)
		}
	default:
		if ok := cm.lockTable.xLock(ctx, k); !ok {
			return fmt.Errorf("unable to XLock %v: %w", k, ErrLockAbort)
		}
	}

	cm.locks[k] = lockTypeX

	return nil
}

func (cm *ConcurrencyMgr[K]) Release() {
	for k := range cm.locks {
		cm.lockTable.release(k)
	}

	clear(cm.locks)
}

func (cm *ConcurrencyMgr[K]) hasAnyLock(k K) bool {
	l := cm.locks[k]
	return l > lockTypeNone
}
