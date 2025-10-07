package async

import "context"

type Handler[T any, U any] chan *Request[T, U]

func (h Handler[T, U]) Send(ctx context.Context, t T) (U, error) {
	r := NewRequest[T, U](ctx, t)
	if !h.send(ctx, r) {
		var u U
		return u, context.Canceled
	}

	return r.wait(ctx)
}

func (h Handler[T, U]) send(ctx context.Context, r *Request[T, U]) bool {
	select {
	case h <- r:
		return true
	case <-ctx.Done():
		return false
	}
}

func (h Handler[T, U]) Listen(ctx context.Context) (*Request[T, U], bool) {
	select {
	case r := <-h:
		return r, true
	case <-ctx.Done():
		return nil, false
	}
}
