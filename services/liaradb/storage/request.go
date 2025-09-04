package storage

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
