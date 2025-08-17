package locktable

import (
	"context"
)

type lockRequestType int8

const (
	lockRequestTypeS lockRequestType = iota
	lockRequestTypeX
	lockRequestTypeUpgrade
)

type lockRequest[K comparable] struct {
	key      K
	_type    lockRequestType
	response chan struct{}
	ctx      context.Context
}

func newExclusiveLockRequest[K comparable](ctx context.Context, k K) *lockRequest[K] {
	return &lockRequest[K]{
		key:      k,
		_type:    lockRequestTypeX,
		response: make(chan struct{}, 1),
		ctx:      ctx,
	}
}

func newUpgradeLockRequest[K comparable](ctx context.Context, k K) *lockRequest[K] {
	return &lockRequest[K]{
		key:      k,
		_type:    lockRequestTypeUpgrade,
		response: make(chan struct{}, 1),
		ctx:      ctx,
	}
}

func newSharedLockRequest[K comparable](ctx context.Context, k K) *lockRequest[K] {
	return &lockRequest[K]{
		key:      k,
		_type:    lockRequestTypeS,
		response: make(chan struct{}, 1),
		ctx:      ctx,
	}
}

func (lr *lockRequest[K]) lockRequestType() lockRequestType {
	return lr._type
}

func (lr *lockRequest[K]) isShared() bool {
	return lr._type == lockRequestTypeS
}

func (lr *lockRequest[K]) wait(ctx context.Context) bool {
	select {
	case <-lr.response:
		lr.response = nil
		return true
	case <-ctx.Done():
		lr.response = nil
		return false
	}
}

func (lr *lockRequest[K]) send() bool {
	select {
	case <-lr.ctx.Done():
		return false
	case lr.response <- struct{}{}:
		return true
	}
}
