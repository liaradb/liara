package controller

import (
	"context"

	"github.com/cardboardrobots/esgrpc"
	"github.com/cardboardrobots/eventsource"
	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara_service/feature/eventsource/service"
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

func (esc *EventSourceController) Append(
	ctx context.Context,
	request *pb.AppendRequest,
) (*pb.AppendResponse, error) {
	events := make([]eventsource.Event, 0, len(request.Events))
	for _, e := range request.Events {
		events = append(events, esgrpc.DtoToEvent(e))
	}
	err := esc.eventService.Append(ctx, events...)
	if err != nil {
		return nil, err
	}

	return &pb.AppendResponse{}, nil
}

func (esc *EventSourceController) Get(
	request *pb.GetRequest,
	stream pb.EventSourceService_GetServer,
) error {
	rows := esc.eventService.Get(stream.Context(),
		eventsource.AggregateID(request.AggregateId))
	for row, err := range rows {
		if err != nil {
			return err
		}

		stream.Send(esgrpc.EventToDto(row))
	}
	return nil
}

func (esc *EventSourceController) GetByAggregateIDAndName(
	request *pb.GetByAggregateIDAndNameRequest,
	stream pb.EventSourceService_GetByAggregateIDAndNameServer,
) error {
	rows := esc.eventService.GetByAggregateIDAndName(stream.Context(),
		eventsource.AggregateID(request.AggregateId),
		eventsource.AggregateName(request.Name))
	for row, err := range rows {
		if err != nil {
			return err
		}

		stream.Send(esgrpc.EventToDto(row))
	}
	return nil
}
