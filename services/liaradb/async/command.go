package async

import "context"

type Command[T any] struct {
	value    T
	response chan response[T, struct{}]
}

func NewCommand[T any](value T) *Command[T] {
	return &Command[T]{
		value:    value,
		response: make(chan response[T, struct{}], 1),
	}
}

func (r *Command[T]) Value() T { return r.value }

func (r *Command[T]) Reply(err error) {
	r.response <- response[T, struct{}]{err: err}
}

func (r *Command[T]) Wait(ctx context.Context) error {
	select {
	case res := <-r.response:
		return res.err
	case <-ctx.Done():
		return context.Canceled
	}
}
