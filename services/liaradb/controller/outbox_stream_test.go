package controller

import (
	"context"

	liara "github.com/liaradb/eventsource_go/generated"
	"google.golang.org/grpc/metadata"
)

type testOutboxStream struct {
	ctx     context.Context
	handler func(*liara.Outbox)
}

var _ liara.EventSourceService_ListOutboxesServer = (*testOutboxStream)(nil)

func newTestOutboxStream(ctx context.Context, h func(*liara.Outbox)) *testOutboxStream {
	return &testOutboxStream{
		ctx:     ctx,
		handler: h,
	}
}

func (os *testOutboxStream) Context() context.Context {
	return os.ctx
}

func (os *testOutboxStream) RecvMsg(m any) error {
	panic("unimplemented")
}

func (os *testOutboxStream) Send(o *liara.Outbox) error {
	os.handler(o)
	return nil
}

func (os *testOutboxStream) SendHeader(metadata.MD) error {
	panic("unimplemented")
}

func (os *testOutboxStream) SendMsg(m any) error {
	panic("unimplemented")
}

func (os *testOutboxStream) SetHeader(metadata.MD) error {
	panic("unimplemented")
}

func (os *testOutboxStream) SetTrailer(metadata.MD) {
	panic("unimplemented")
}
