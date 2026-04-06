package async

import (
	"context"
	"errors"
)

type Request[T any, U any] struct {
	ctx      context.Context
	value    T
	response chan response[T, U]
	onCancel func(U, error) error
}

type response[T any, U any] struct {
	value U
	err   error
}

func newRequest[T any, U any](ctx context.Context, value T, onCancel func(U, error) error) *Request[T, U] {
	return &Request[T, U]{
		ctx:      ctx,
		value:    value,
		response: make(chan response[T, U], 1),
		onCancel: onCancel,
	}
}

func (r *Request[T, U]) Context() context.Context { return r.ctx }
func (r *Request[T, U]) Value() T                 { return r.value }

func (r *Request[T, U]) Reply(value U, err error) error {
	select {
	case r.response <- response[T, U]{
		value: value,
		err:   err}:
		return nil
	case <-r.ctx.Done():
		return errors.Join(context.Canceled, r.onCancel(value, err))
	}
}

func (r *Request[T, U]) wait(ctx context.Context) (U, error) {
	select {
	case res := <-r.response:
		return res.value, res.err
	case <-ctx.Done():
		var v U

		if len(r.response) > 0 && r.onCancel != nil {
			res := <-r.response
			err := r.onCancel(res.value, res.err)
			return v, errors.Join(context.Canceled, err)
		}

		return v, context.Canceled
	}
}
