package controller

import (
	"context"

	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/service"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
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
	err := esc.eventService.Append(ctx,
		mapSlice(request.Events, dtoToAppendEvent)...)
	if err != nil {
		return nil, err
	}

	return &pb.AppendResponse{}, nil
}

func (esc *EventSourceController) Get(
	request *pb.GetRequest,
	stream pb.EventSourceService_GetServer,
) error {
	for row, err := range esc.eventService.Get(stream.Context(),
		value.AggregateID(request.AggregateId)) {
		if err != nil {
			return err
		}

		stream.Send(eventToDto(row))
	}
	return nil
}

func (esc *EventSourceController) GetByAggregateIDAndName(
	request *pb.GetByAggregateIDAndNameRequest,
	stream pb.EventSourceService_GetByAggregateIDAndNameServer,
) error {
	for row, err := range esc.eventService.GetByAggregateIDAndName(stream.Context(),
		value.AggregateID(request.AggregateId),
		value.AggregateName(request.Name)) {
		if err != nil {
			return err
		}

		stream.Send(eventToDto(row))
	}
	return nil
}

func (esc *EventSourceController) GetAfterGlobalVersion(
	request *pb.GetAfterGlobalVersionRequest,
	stream pb.EventSourceService_GetAfterGlobalVersionServer,
) error {
	for row, err := range esc.eventService.GetAfterGlobalVersion(stream.Context(),
		value.GlobalVersion(request.GlobalVersion),
		value.PartitionID(request.PartitionId),
		value.Limit(request.Limit)) {
		if err != nil {
			return err
		}

		stream.Send(eventToDto(row))
	}
	return nil
}

func (esc *EventSourceController) GetOrCreateOutbox(
	ctx context.Context,
	request *pb.GetOrCreateOutboxRequest,
) (*pb.GetOrCreateOutboxResponse, error) {
	result, err := esc.eventService.GetOrCreateOutbox(ctx, value.OutboxID(request.OutboxId))
	if err != nil {
		return nil, err
	}

	return &pb.GetOrCreateOutboxResponse{
		GlobalVersion: int64(result),
	}, nil
}

func (esc *EventSourceController) UpdateOutboxPosition(
	ctx context.Context,
	request *pb.UpdateOutboxPositionRequest,
) (*pb.UpdateOutboxPositionResponse, error) {
	err := esc.eventService.UpdateOutboxPosition(ctx,
		value.OutboxID(request.OutboxId),
		value.GlobalVersion(request.GlobalVersion))
	if err != nil {
		return nil, err
	}

	return &pb.UpdateOutboxPositionResponse{}, nil
}

func mapSlice[T any, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, 0, len(slice))
	for _, item := range slice {
		result = append(result, mapper(item))
	}
	return result
}
