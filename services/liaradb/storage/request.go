package storage

import "context"

type request struct {
	ctx      context.Context
	value    BlockID
	response chan *response
}

type response struct {
	value *Buffer
	err   error
}

// External thread
func newRequest(ctx context.Context, bid BlockID) *request {
	return &request{
		ctx:      ctx,
		value:    bid,
		response: make(chan *response, 1), // TODO: Test this async
	}
}

func (r *request) respond(b *Buffer, err error) {
	select {
	case <-r.ctx.Done():
	case r.response <- &response{
		value: b,
		err:   err}:
	}
}

// External thread
func (r *request) wait(ctx context.Context) (*Buffer, error) {
	select {
	case res := <-r.response:
		return res.value, res.err
	case <-ctx.Done():
		return nil, context.Canceled
	}
}
