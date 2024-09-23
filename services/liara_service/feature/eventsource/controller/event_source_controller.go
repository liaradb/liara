package controller

import (
	"context"

	pb "github.com/cardboardrobots/eventsource_go/generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventSourceController struct {
	pb.UnimplementedEventSourceServiceServer
}

func (esc *EventSourceController) ListEvents(ctx context.Context, request *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEvents not implemented")
}
