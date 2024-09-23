package esgrpc

import (
	"context"
	"io"
	"iter"

	"github.com/cardboardrobots/eventsource"
	pb "github.com/cardboardrobots/eventsource_go/generated"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventSourceGRPC struct {
	client pb.EventSourceServiceClient
}

var _ eventsource.EventSource = &EventSourceGRPC{}

func NewEventSourceGRPC(
	conn *grpc.ClientConn,
) *EventSourceGRPC {
	client := pb.NewEventSourceServiceClient(conn)

	return &EventSourceGRPC{
		client: client,
	}
}

func (es *EventSourceGRPC) Append(ctx context.Context, events ...eventsource.Event) error {
	data := make([]*pb.Event, 0, len(events))
	for _, event := range events {
		data = append(data, EventToDto(event))
	}

	_, err := es.client.Append(ctx, &pb.AppendRequest{
		Events: data,
	})
	return err
}

func (es *EventSourceGRPC) Get(ctx context.Context, id eventsource.AggregateID) iter.Seq2[eventsource.Event, error] {
	return func(yield func(eventsource.Event, error) bool) {
		stream, err := es.client.Get(ctx, &pb.GetRequest{
			AggregateId: string(id),
		})
		if err != nil {
			yield(eventsource.Event{}, err)
			return
		}

		for {
			e, err := stream.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				yield(eventsource.Event{}, err)
				return
			}

			yield(DtoToEvent(e), nil)
		}
	}
}

func (es *EventSourceGRPC) GetByAggregateIDAndName(ctx context.Context, id eventsource.AggregateID, name eventsource.AggregateName) iter.Seq2[eventsource.Event, error] {
	return func(yield func(eventsource.Event, error) bool) {
		stream, err := es.client.GetByAggregateIDAndName(ctx, &pb.GetByAggregateIDAndNameRequest{
			AggregateId: string(id),
			Name:        string(name),
		})
		if err != nil {
			yield(eventsource.Event{}, err)
			return
		}

		for {
			e, err := stream.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				yield(eventsource.Event{}, err)
				return
			}

			yield(DtoToEvent(e), nil)
		}
	}
}

func DtoToEvent(dto *pb.Event) eventsource.Event {
	return eventsource.Event{
		GlobalVersion: eventsource.GlobalVersion(dto.GlobalVersion),
		AggregateName: eventsource.AggregateName(dto.AggregateName),
		ID:            eventsource.EventID(dto.Id),
		AggregateID:   eventsource.AggregateID(dto.AggregateId),
		Version:       eventsource.Version(dto.Version),
		Name:          eventsource.EventName(dto.Name),
		CorrelationID: eventsource.CorrelationID(dto.CorrelationId),
		UserID:        eventsource.UserID(dto.UserId),
		Time:          dto.Time.AsTime(),
		Schema:        eventsource.Schema(dto.Schema),
		Data:          dto.Data,
	}
}

func EventToDto(e eventsource.Event) *pb.Event {
	return &pb.Event{
		GlobalVersion: int64(e.GlobalVersion),
		AggregateName: string(e.AggregateName),
		Id:            string(e.ID),
		AggregateId:   string(e.AggregateID),
		Version:       int64(e.Version),
		Name:          e.Name.String(),
		CorrelationId: string(e.CorrelationID),
		UserId:        string(e.UserID),
		Time:          timestamppb.New(e.Time),
		Schema:        string(e.Schema),
		Data:          e.Data,
	}
}
