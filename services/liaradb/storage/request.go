package storage

import "context"

type request struct {
	blockID BlockID
	out     chan *response
}

type response struct {
	buffer *Buffer
	err    error
}

func newRequest(bid BlockID) *request {
	return &request{
		blockID: bid,
		out:     make(chan *response), // TODO: Make this async
	}
}

func (r *request) close() {
	close(r.out)
}

func (r *request) respond(b *Buffer, err error) {
	r.out <- &response{
		buffer: b,
		err:    err,
	}
}

func (r *request) wait(ctx context.Context) (*Buffer, error) {
	select {
	case o, ok := <-r.out:
		if ok {
			return o.buffer, o.err
		} else {
			return nil, ErrRequestClosed
		}
	case <-ctx.Done():
		return nil, context.Canceled
	}
}
