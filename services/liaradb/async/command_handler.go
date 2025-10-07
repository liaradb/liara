package async

import "context"

type CommandHandler[T any] chan *Command[T]

func (h CommandHandler[T]) Send(ctx context.Context, t T) error {
	r := NewCommand(ctx, t)
	if !h.send(ctx, r) {
		return context.Canceled
	}

	return r.wait(ctx)
}

func (h CommandHandler[T]) send(ctx context.Context, c *Command[T]) bool {
	select {
	case h <- c:
		return true
	case <-ctx.Done():
		return false
	}
}
