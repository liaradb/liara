package esgrpc

import (
	"context"
	"io"
	"iter"

	"github.com/cardboardrobots/eventsource/entity"
	"github.com/cardboardrobots/eventsource/service"
	"github.com/cardboardrobots/eventsource/value"
	pb "github.com/cardboardrobots/eventsource_go/generated"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventSourceGRPC struct {
	client pb.EventSourceServiceClient
}

var _ service.EventSource = &EventSourceGRPC{}
var _ service.EventRepository = &EventSourceGRPC{}
var _ service.OutboxRepository = &EventSourceGRPC{}

func NewEventSourceGRPC(
	conn *grpc.ClientConn,
) *EventSourceGRPC {
	client := pb.NewEventSourceServiceClient(conn)

	return &EventSourceGRPC{
		client: client,
	}
}

func (es *EventSourceGRPC) Append(
	ctx context.Context,
	events ...entity.AppendEvent,
) error {
	data := make([]*pb.AppendEvent, 0, len(events))
	for _, event := range events {
		data = append(data, AppendEventToDto(event))
	}

	_, err := es.client.Append(ctx, &pb.AppendRequest{
		Events: data,
	})
	return err
}

func (es *EventSourceGRPC) Get(
	ctx context.Context,
	id value.AggregateID,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {
		stream, err := es.client.Get(ctx, &pb.GetRequest{
			AggregateId: id.String(),
		})
		if err != nil {
			yield(entity.Event{}, err)
			return
		}

		for {
			e, err := stream.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				yield(entity.Event{}, err)
				return
			}

			if !yield(DtoToEvent(e), nil) {
				return
			}
		}
	}
}

func (es *EventSourceGRPC) GetByAggregateIDAndName(
	ctx context.Context,
	id value.AggregateID,
	name value.AggregateName,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {
		stream, err := es.client.GetByAggregateIDAndName(ctx, &pb.GetByAggregateIDAndNameRequest{
			AggregateId: id.String(),
			Name:        name.String(),
		})
		if err != nil {
			yield(entity.Event{}, err)
			return
		}

		for {
			e, err := stream.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				yield(entity.Event{}, err)
				return
			}

			if !yield(DtoToEvent(e), nil) {
				return
			}
		}
	}
}

func (es *EventSourceGRPC) GetAfterGlobalVersion(
	ctx context.Context,
	version value.GlobalVersion,
	limit value.Limit,
) iter.Seq2[entity.Event, error] {
	return func(yield func(entity.Event, error) bool) {
		stream, err := es.client.GetAfterGlobalVersion(ctx, &pb.GetAfterGlobalVersionRequest{
			GlobalVersion: int64(version),
			Limit:         int64(limit),
		})
		if err != nil {
			yield(entity.Event{}, err)
			return
		}

		for {
			e, err := stream.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				yield(entity.Event{}, err)
				return
			}

			if !yield(DtoToEvent(e), nil) {
				return
			}
		}
	}
}

func (es *EventSourceGRPC) GetOrCreateOutbox(
	ctx context.Context,
	outboxID value.OutboxID,
) (value.GlobalVersion, error) {
	response, err := es.client.GetOrCreateOutbox(ctx, &pb.GetOrCreateOutboxRequest{
		OutboxId: outboxID.String(),
	})
	if err != nil {
		return 0, err
	}

	return value.GlobalVersion(response.GlobalVersion), nil
}

func (es *EventSourceGRPC) UpdateOutboxPosition(
	ctx context.Context,
	outboxID value.OutboxID,
	globalVersion value.GlobalVersion,
) error {
	_, err := es.client.UpdateOutboxPosition(ctx, &pb.UpdateOutboxPositionRequest{
		OutboxId:      outboxID.String(),
		GlobalVersion: int64(globalVersion),
	})
	return err
}

func DtoToEvent(dto *pb.Event) entity.Event {
	return entity.Event{
		GlobalVersion: value.GlobalVersion(dto.GlobalVersion),
		AggregateName: value.AggregateName(dto.AggregateName),
		ID:            value.EventID(dto.Id),
		AggregateID:   value.AggregateID(dto.AggregateId),
		Version:       value.Version(dto.Version),
		Name:          value.EventName(dto.Name),
		CorrelationID: value.CorrelationID(dto.CorrelationId),
		UserID:        value.UserID(dto.UserId),
		Time:          dto.Time.AsTime(),
		Schema:        value.Schema(dto.Schema),
		Data:          dto.Data,
	}
}

func EventToDto(e entity.Event) *pb.Event {
	return &pb.Event{
		GlobalVersion: int64(e.GlobalVersion),
		AggregateName: e.AggregateName.String(),
		Id:            e.ID.String(),
		AggregateId:   e.AggregateID.String(),
		Version:       int64(e.Version),
		Name:          e.Name.String(),
		CorrelationId: e.CorrelationID.String(),
		UserId:        e.UserID.String(),
		Time:          timestamppb.New(e.Time),
		Schema:        e.Schema.String(),
		Data:          e.Data,
	}
}

func DtoToAppendEvent(dto *pb.AppendEvent) entity.AppendEvent {
	return entity.AppendEvent{
		AggregateName: value.AggregateName(dto.AggregateName),
		ID:            value.EventID(dto.Id),
		AggregateID:   value.AggregateID(dto.AggregateId),
		Version:       value.Version(dto.Version),
		Name:          value.EventName(dto.Name),
		CorrelationID: value.CorrelationID(dto.CorrelationId),
		UserID:        value.UserID(dto.UserId),
		Time:          dto.Time.AsTime(),
		Schema:        value.Schema(dto.Schema),
		Data:          dto.Data,
	}
}

func AppendEventToDto(e entity.AppendEvent) *pb.AppendEvent {
	return &pb.AppendEvent{
		AggregateName: e.AggregateName.String(),
		Id:            e.ID.String(),
		AggregateId:   e.AggregateID.String(),
		Version:       int64(e.Version),
		Name:          e.Name.String(),
		CorrelationId: e.CorrelationID.String(),
		UserId:        e.UserID.String(),
		Time:          timestamppb.New(e.Time),
		Schema:        e.Schema.String(),
		Data:          e.Data,
	}
}
