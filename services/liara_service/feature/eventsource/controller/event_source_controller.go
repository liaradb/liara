package controller

import (
	"context"

	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara_service/feature/eventsource/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventSourceController struct {
	pb.UnimplementedEventSourceServiceServer
	eventService *service.EventService
}

func NewEventSourceController(
	eventService *service.EventService,
) *EventSourceController {
	return &EventSourceController{
		eventService: eventService,
	}
}

func (esc *EventSourceController) Append(ctx context.Context, stream *pb.AppendRequest) (*pb.AppendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}

func (esc *EventSourceController) Get(request *pb.GetRequest, stream pb.EventSourceService_GetServer) error {
	_ = esc.eventService.Get(stream.Context())
	return status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func (esc *EventSourceController) GetByAggregateIDAndName(request *pb.GetByAggregateIDAndNameRequest, stream pb.EventSourceService_GetByAggregateIDAndNameServer) error {
	return status.Errorf(codes.Unimplemented, "method GetByAggregateIDAndName not implemented")
}
