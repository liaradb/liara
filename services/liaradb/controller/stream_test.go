package controller

import (
	"context"

	liara "github.com/liaradb/eventsource_go/generated"
	"google.golang.org/grpc/metadata"
)

type testStream[T any] struct {
	ctx     context.Context
	handler func(T)
}

var _ liara.EventSourceService_ListOutboxesServer = (*testStream[*liara.Outbox])(nil)

func newTestStream[T any](ctx context.Context, h func(T)) *testStream[T] {
	return &testStream[T]{
		ctx:     ctx,
		handler: h,
	}
}

func (os *testStream[T]) Send(t T) error {
	os.handler(t)
	return nil
}

func (os *testStream[T]) Context() context.Context {
	return os.ctx
}

func (os *testStream[T]) RecvMsg(m any) error {
	panic("unimplemented")
}

func (os *testStream[T]) SendHeader(metadata.MD) error {
	panic("unimplemented")
}

func (os *testStream[T]) SendMsg(m any) error {
	panic("unimplemented")
}

func (os *testStream[T]) SetHeader(metadata.MD) error {
	panic("unimplemented")
}

func (os *testStream[T]) SetTrailer(metadata.MD) {
	panic("unimplemented")
}
