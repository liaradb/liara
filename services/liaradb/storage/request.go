package storage

import "context"

type request struct {
	blockID BlockID
	ctx     context.Context
	out     chan *response
}

type response struct {
	buffer *Buffer
	err    error
}

// External thread
func newRequest(ctx context.Context, bid BlockID) *request {
	return &request{
		blockID: bid,
		ctx:     ctx,
		out:     make(chan *response), // TODO: Make this async
	}
}

func (r *request) respond(ctx context.Context, b *Buffer, err error) {
	select {
	case r.out <- &response{
		buffer: b,
		err:    err,
	}:
	case <-ctx.Done():
	}
}

// External thread
func (r *request) wait(ctx context.Context) (*Buffer, error) {
	select {
	case o := <-r.out:
		return o.buffer, o.err
	case <-ctx.Done():
		return nil, context.Canceled
	}
}
