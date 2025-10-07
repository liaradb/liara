package async

import "context"

type Command[T any] struct {
	ctx      context.Context
	value    T
	response chan response[T, struct{}]
}

func NewCommand[T any](ctx context.Context, value T) *Command[T] {
	return &Command[T]{
		ctx:      ctx,
		value:    value,
		response: make(chan response[T, struct{}], 1),
	}
}

func (r *Command[T]) Value() T { return r.value }

func (r *Command[T]) Reply(err error) {
	select {
	case r.response <- response[T, struct{}]{err: err}:
	case <-r.ctx.Done():
	}
}

func (r *Command[T]) wait(ctx context.Context) error {
	select {
	case res := <-r.response:
		return res.err
	case <-ctx.Done():
		return context.Canceled
	}
}
