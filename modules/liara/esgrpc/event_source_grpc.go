package esgrpc

import (
	"context"
	"io"
	"iter"

	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara"
	"google.golang.org/grpc"
)

type EventSourceGRPC struct {
	client pb.EventSourceServiceClient
}

var _ liara.EventRepository = &EventSourceGRPC{}
var _ liara.OutboxRepository = &EventSourceGRPC{}

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
	events ...liara.AppendEvent,
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
	id liara.AggregateID,
) iter.Seq2[liara.Event, error] {
	return func(yield func(liara.Event, error) bool) {
		stream, err := es.client.Get(ctx, &pb.GetRequest{
			AggregateId: id.String(),
		})
		if err != nil {
			yield(liara.Event{}, err)
			return
		}

		for {
			e, err := stream.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				yield(liara.Event{}, err)
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
	id liara.AggregateID,
	name liara.AggregateName,
) iter.Seq2[liara.Event, error] {
	return func(yield func(liara.Event, error) bool) {
		stream, err := es.client.GetByAggregateIDAndName(ctx, &pb.GetByAggregateIDAndNameRequest{
			AggregateId: id.String(),
			Name:        name.String(),
		})
		if err != nil {
			yield(liara.Event{}, err)
			return
		}

		for {
			e, err := stream.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				yield(liara.Event{}, err)
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
	version liara.GlobalVersion,
	partitionIDs []liara.PartitionID,
	limit liara.Limit,
) iter.Seq2[liara.Event, error] {
	return func(yield func(liara.Event, error) bool) {
		pids := make([]int32, 0, len(partitionIDs))
		for _, p := range partitionIDs {
			pids = append(pids, p.Value())
		}
		stream, err := es.client.GetAfterGlobalVersion(ctx, &pb.GetAfterGlobalVersionRequest{
			GlobalVersion: int64(version),
			PartitionIds:  pids,
			Limit:         int64(limit),
		})
		if err != nil {
			yield(liara.Event{}, err)
			return
		}

		for {
			e, err := stream.Recv()
			if err == io.EOF {
				return
			}

			if err != nil {
				yield(liara.Event{}, err)
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
	outboxID liara.OutboxID,
	partitionIDs []liara.PartitionID,
) (liara.GlobalVersion, error) {
	response, err := es.client.GetOrCreateOutbox(ctx, &pb.GetOrCreateOutboxRequest{
		OutboxId: outboxID.String(),
	})
	if err != nil {
		return 0, err
	}

	return liara.GlobalVersion(response.GlobalVersion), nil
}

func (es *EventSourceGRPC) UpdateOutboxPosition(
	ctx context.Context,
	outboxID liara.OutboxID,
	globalVersion liara.GlobalVersion,
) error {
	_, err := es.client.UpdateOutboxPosition(ctx, &pb.UpdateOutboxPositionRequest{
		OutboxId:      outboxID.String(),
		GlobalVersion: int64(globalVersion),
	})
	return err
}
