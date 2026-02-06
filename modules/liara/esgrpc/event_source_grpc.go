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
	tenantID liara.TenantID,
	options liara.AppendOptions,
	events ...liara.AppendEvent,
) error {
	data := make([]*pb.AppendEvent, 0, len(events))
	for _, event := range events {
		data = append(data, appendEventToDto(event))
	}

	_, err := es.client.Append(ctx, &pb.AppendRequest{
		TenantId: tenantID.String(),
		Options:  appendOptionsToDto(options),
		Events:   data,
	})
	return err
}

func (es *EventSourceGRPC) TestIdempotency(
	ctx context.Context,
	tenantID liara.TenantID,
	requetID liara.RequestID,
) (bool, error) {
	response, err := es.client.TestIdempotency(ctx, &pb.TestIdempotencyRequest{
		TenantId:  requetID.String(),
		RequestId: requetID.String(),
	})
	return response.Ok, err
}

func (es *EventSourceGRPC) Get(
	ctx context.Context,
	tenantID liara.TenantID,
	id liara.AggregateID,
) iter.Seq2[liara.Event, error] {
	return func(yield func(liara.Event, error) bool) {
		stream, err := es.client.Get(ctx, &pb.GetRequest{
			TenantId:    tenantID.String(),
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
	tenantID liara.TenantID,
	id liara.AggregateID,
	name liara.AggregateName,
) iter.Seq2[liara.Event, error] {
	return func(yield func(liara.Event, error) bool) {
		stream, err := es.client.GetByAggregateIDAndName(ctx, &pb.GetByAggregateIDAndNameRequest{
			TenantId:    tenantID.String(),
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
	tenantID liara.TenantID,
	version liara.GlobalVersion,
	low liara.PartitionID,
	high liara.PartitionID,
	limit liara.Limit,
) iter.Seq2[liara.Event, error] {
	return func(yield func(liara.Event, error) bool) {
		stream, err := es.client.GetAfterGlobalVersion(ctx, &pb.GetAfterGlobalVersionRequest{
			TenantId:      tenantID.String(),
			GlobalVersion: int64(version),
			Low:           low.Value(),
			High:          high.Value(),
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

func (es *EventSourceGRPC) GetByOutbox(
	ctx context.Context,
	tenantID liara.TenantID,
	outboxID liara.OutboxID,
	limit liara.Limit,
) iter.Seq2[liara.Event, error] {
	return func(yield func(liara.Event, error) bool) {
		stream, err := es.client.GetByOutbox(ctx, &pb.GetByOutboxRequest{
			TenantId: tenantID.String(),
			OutboxId: outboxID.String(),
			Limit:    int64(limit),
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

func (es *EventSourceGRPC) CreateOutbox(
	ctx context.Context,
	tenantID liara.TenantID,
	low liara.PartitionID,
	high liara.PartitionID,
) (liara.OutboxID, error) {
	response, err := es.client.CreateOutbox(ctx, &pb.CreateOutboxRequest{
		TenantId: tenantID.String(),
		Low:      low.Value(),
		High:     high.Value(),
	})
	if err != nil {
		return "", err
	}

	return liara.OutboxID(response.OutboxId), nil
}

func (es *EventSourceGRPC) GetOutbox(
	ctx context.Context,
	tenantID liara.TenantID,
	outboxID liara.OutboxID,
) (liara.GlobalVersion, error) {
	response, err := es.client.GetOutbox(ctx, &pb.GetOutboxRequest{
		TenantId: tenantID.String(),
		OutboxId: outboxID.String(),
	})
	if err != nil {
		return 0, err
	}

	return liara.GlobalVersion(response.GlobalVersion), nil
}

func (es *EventSourceGRPC) UpdateOutboxPosition(
	ctx context.Context,
	tenantID liara.TenantID,
	outboxID liara.OutboxID,
	globalVersion liara.GlobalVersion,
) error {
	_, err := es.client.UpdateOutboxPosition(ctx, &pb.UpdateOutboxPositionRequest{
		TenantId:      tenantID.String(),
		OutboxId:      outboxID.String(),
		GlobalVersion: int64(globalVersion),
	})
	return err
}

func (es *EventSourceGRPC) CreateTenant(
	ctx context.Context,
	tenantName liara.TenantName,
) (liara.TenantID, error) {
	response, err := es.client.CreateTenant(ctx, &pb.CreateTenantRequest{
		Name: tenantName.String(),
	})
	if err != nil {
		return "", err
	}

	return liara.TenantID(response.TenantId), nil
}

func (es *EventSourceGRPC) DeleteTenant(
	ctx context.Context,
	tenantID liara.TenantID,
) error {
	_, err := es.client.DeleteTenant(ctx, &pb.DeleteTenantRequest{
		TenantId: tenantID.String(),
	})
	return err
}

func (es *EventSourceGRPC) ListTenants(
	ctx context.Context,
) iter.Seq2[*pb.Tenant, error] {
	return func(yield func(*pb.Tenant, error) bool) {
		result, err := es.client.ListTenants(ctx, &pb.ListTenantsRequest{})
		if err != nil {
			yield(nil, err)
			return
		}

		for {
			m, err := result.Recv()
			if err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				return
			}

			if !yield(m, nil) {
				return
			}
		}
	}
}
