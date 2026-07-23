package async

import "context"

type Command[T any] struct {
	value    T
	response chan response[T, struct{}]
}

func newCommand[T any](value T) *Command[T] {
	return &Command[T]{
		value:    value,
		response: make(chan response[T, struct{}], 1),
	}
}

func (r *Command[T]) Value() T { return r.value }

func (r *Command[T]) Reply(err error) {
	r.response <- response[T, struct{}]{err: err}
	close(r.response)
}

func (r *Command[T]) wait(ctx context.Context) error {
	select {
	case res := <-r.response:
		return res.err
	case <-ctx.Done():
		return context.Canceled
	}
}
