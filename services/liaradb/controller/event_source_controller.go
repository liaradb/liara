package controller

import (
	"context"

	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
)

type EventSourceController struct {
	pb.UnimplementedEventSourceServiceServer
	eventService  *service.EventService
	tenantService *service.TenantService
}

func NewEventSourceController(
	eventService *service.EventService,
	tenantService *service.TenantService,
) *EventSourceController {
	return &EventSourceController{
		eventService:  eventService,
		tenantService: tenantService,
	}
}

func (esc *EventSourceController) Append(
	ctx context.Context,
	request *pb.AppendRequest,
) (*pb.AppendResponse, error) {
	err := esc.eventService.Append(ctx,
		value.TenantID(request.TenantId),
		dtoToAppendOptions(request.Options),
		value.NewPartitionID(0),
		mapSlice(request.Events, dtoToAppendEvent)...)
	if err != nil {
		return nil, err
	}

	return &pb.AppendResponse{}, nil
}

func (esc *EventSourceController) TestIdempotency(
	ctx context.Context,
	request *pb.TestIdempotencyRequest,
) (*pb.TestIdempotencyResponse, error) {
	ok, err := esc.eventService.TestIdempotency(ctx,
		value.TenantID(request.TenantId),
		value.RequestID(request.RequestId))
	if err != nil {
		return nil, err
	}

	return &pb.TestIdempotencyResponse{
		Ok: ok,
	}, nil
}

func (esc *EventSourceController) Get(
	request *pb.GetRequest,
	stream pb.EventSourceService_GetServer,
) error {
	for row, err := range esc.eventService.Get(stream.Context(),
		value.TenantID(request.TenantId),
		value.NewPartitionID(0), // TODO: Change this to parameter
		value.NewAggregateID(request.AggregateId),
	) {
		if err != nil {
			return err
		}

		if err := stream.Send(eventToDto(row)); err != nil {
			return err
		}
	}
	return nil
}

func (esc *EventSourceController) GetByAggregateIDAndName(
	request *pb.GetByAggregateIDAndNameRequest,
	stream pb.EventSourceService_GetByAggregateIDAndNameServer,
) error {
	for row, err := range esc.eventService.GetByAggregateIDAndName(stream.Context(),
		value.TenantID(request.TenantId),
		value.NewPartitionID(0), // TODO: Change this to parameter
		value.NewAggregateID(request.AggregateId),
		value.NewAggregateName(request.Name)) {
		if err != nil {
			return err
		}

		if err := stream.Send(eventToDto(row)); err != nil {
			return err
		}
	}
	return nil
}

func (esc *EventSourceController) GetAfterGlobalVersion(
	request *pb.GetAfterGlobalVersionRequest,
	stream pb.EventSourceService_GetAfterGlobalVersionServer,
) error {
	for row, err := range esc.eventService.GetAfterGlobalVersion(stream.Context(),
		value.TenantID(request.TenantId),
		value.NewGlobalVersion(uint64(request.GlobalVersion)),
		dtoToPartitionRange(request.PartitionIds),
		value.Limit(request.Limit)) {
		if err != nil {
			return err
		}

		if err := stream.Send(eventToDto(row)); err != nil {
			return err
		}
	}
	return nil
}

func (esc *EventSourceController) GetByOutbox(
	request *pb.GetByOutboxRequest,
	stream pb.EventSourceService_GetByOutboxServer,
) error {
	oid, err := value.NewOutboxIDFromString(request.OutboxId)
	if err != nil {
		return err
	}

	// TODO: How do we specify partition?
	pid := value.NewPartitionID(0)
	for row, err := range esc.eventService.GetByOutbox(stream.Context(),
		value.TenantID(request.TenantId),
		pid,
		oid,
		value.Limit(request.Limit)) {
		if err != nil {
			return err
		}

		if err := stream.Send(eventToDto(row)); err != nil {
			return err
		}
	}
	return nil
}

func (esc *EventSourceController) CreateOutbox(
	ctx context.Context,
	request *pb.CreateOutboxRequest,
) (*pb.CreateOutboxResponse, error) {
	oid, err := value.NewOutboxIDFromString(request.OutboxId)
	if err != nil {
		return nil, err
	}

	outboxID, err := esc.eventService.CreateOutbox(ctx,
		value.TenantID(request.TenantId),
		oid,
		dtoToPartitionRange(request.PartitionId))
	if err != nil {
		return nil, err
	}

	return &pb.CreateOutboxResponse{
		OutboxId: outboxID.String(),
	}, nil
}

func (esc *EventSourceController) GetOutbox(
	ctx context.Context,
	request *pb.GetOutboxRequest,
) (*pb.GetOutboxResponse, error) {
	oid, err := value.NewOutboxIDFromString(request.OutboxId)
	if err != nil {
		return nil, err
	}

	// TODO: How do we specify partition?
	pid := value.NewPartitionID(0)

	result, err := esc.eventService.GetOutbox(ctx,
		value.TenantID(request.TenantId),
		pid,
		oid)
	if err != nil {
		return nil, err
	}
	low, high := result.PartitionRange().All()

	return &pb.GetOutboxResponse{
		GlobalVersion: int64(result.GlobalVersion().Value()),
		PartitionId:   []int32{int32(low.Value()), int32(high.Value())},
	}, nil
}

func (esc *EventSourceController) UpdateOutboxPosition(
	ctx context.Context,
	request *pb.UpdateOutboxPositionRequest,
) (*pb.UpdateOutboxPositionResponse, error) {
	oid, err := value.NewOutboxIDFromString(request.OutboxId)
	if err != nil {
		return nil, err
	}

	// TODO: How do we specify partition?
	pid := value.NewPartitionID(0)

	if err := esc.eventService.UpdateOutboxPosition(ctx,
		value.TenantID(request.TenantId),
		pid,
		oid,
		value.NewGlobalVersion(uint64(request.GlobalVersion))); err != nil {
		return nil, err
	}

	return &pb.UpdateOutboxPositionResponse{}, nil
}

func (esc *EventSourceController) CreateTenant(ctx context.Context, request *pb.CreateTenantRequest) (*pb.CreateTenantReponse, error) {
	id, err := esc.tenantService.Create(ctx, service.CreateTenantCommand{
		TenantID:   value.TenantID(request.TenantId),
		TenantName: value.TenantName(request.Name),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateTenantReponse{
		TenantId: id.String(),
	}, nil
}

func (esc *EventSourceController) DeleteTenant(ctx context.Context, request *pb.DeleteTenantRequest) (*pb.DeleteTenantResponse, error) {
	err := esc.tenantService.Delete(ctx, service.DeleteTenantCommand{
		TenantID: value.TenantID(request.TenantId),
	})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteTenantResponse{}, nil
}

func (esc *EventSourceController) RenameTenant(ctx context.Context, request *pb.RenameTenantRequest) (*pb.RenameTenantResponse, error) {
	err := esc.tenantService.Rename(ctx, service.RenameTenantCommand{
		TenantName: value.TenantName(request.Name),
	})
	if err != nil {
		return nil, err
	}

	return &pb.RenameTenantResponse{}, nil
}

func (esc *EventSourceController) GetTenant(ctx context.Context, request *pb.GetTenantRequest) (*pb.GetTenantResponse, error) {
	t, err := esc.tenantService.Get(ctx, value.TenantID(request.TenantId))
	if err != nil {
		return nil, err
	}

	return &pb.GetTenantResponse{
		Tenant: &pb.Tenant{
			TenantId: t.ID().String(),
			Name:     t.Name().String(),
		},
	}, nil
}

func (esc *EventSourceController) ListTenants(request *pb.ListTenantsRequest, stream pb.EventSourceService_ListTenantsServer) error {
	for t, err := range esc.tenantService.List(stream.Context(), 0, 0) {
		if err != nil {
			return err
		}

		stream.Send(&pb.Tenant{
			TenantId: t.ID().String(),
			Name:     t.Name().String(),
		})
	}
	return nil
}

func mapSlice[T any, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, 0, len(slice))
	for _, item := range slice {
		result = append(result, mapper(item))
	}
	return result
}
