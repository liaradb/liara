package async

import "context"

type Request[T any, U any] struct {
	ctx      context.Context
	value    T
	response chan response[T, U]
}

type response[T any, U any] struct {
	value U
	err   error
}

func NewRequest[T any, U any](ctx context.Context, value T) *Request[T, U] {
	return &Request[T, U]{
		ctx:      ctx,
		value:    value,
		response: make(chan response[T, U], 1),
	}
}

func (r *Request[T, U]) Context() context.Context { return r.ctx }
func (r *Request[T, U]) Value() T                 { return r.value }

func (r *Request[T, U]) Reply(value U, err error) {
	select {
	case r.response <- response[T, U]{
		value: value,
		err:   err}:
	case <-r.ctx.Done():
	}
}

func (r *Request[T, U]) Wait(ctx context.Context) (U, error) {
	select {
	case res := <-r.response:
		return res.value, res.err
	case <-ctx.Done():
		var v U
		return v, context.Canceled
	}
}
