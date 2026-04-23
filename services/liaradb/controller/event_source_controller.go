package controller

import (
	"context"

	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
)

type EventSourceController struct {
	pb.UnimplementedEventSourceServiceServer
	eventService  EventService
	tenantService TenantService
}

func NewEventSourceController(
	eventService EventService,
	tenantService TenantService,
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
	o, err := dtoToAppendOptions(request.Options)
	if err != nil {
		return nil, err
	}

	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return nil, err
	}

	if err := esc.eventService.Append(ctx,
		tid,
		o,
		value.NewPartitionID(request.PartitionId),
		mapSlice(request.Events, dtoToAppendEvent)...); err != nil {
		return nil, err
	}

	return &pb.AppendResponse{}, nil
}

func (esc *EventSourceController) TestIdempotency(
	ctx context.Context,
	request *pb.TestIdempotencyRequest,
) (*pb.TestIdempotencyResponse, error) {
	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return nil, err
	}

	rid, err := value.NewRequestIDFromString(request.RequestId)
	if err != nil {
		return nil, err
	}

	ok, err := esc.eventService.TestIdempotency(ctx,
		tid,
		rid)
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
	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return err
	}

	for row, err := range esc.eventService.Get(stream.Context(),
		tid,
		value.NewPartitionID(request.PartitionId),
		value.NewAggregateID(request.AggregateId),
	) {
		if err != nil {
			return err
		}

		if err := stream.Send(eventToDTO(row)); err != nil {
			return err
		}
	}
	return nil
}

func (esc *EventSourceController) GetByAggregateIDAndName(
	request *pb.GetByAggregateIDAndNameRequest,
	stream pb.EventSourceService_GetByAggregateIDAndNameServer,
) error {
	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return err
	}

	for row, err := range esc.eventService.GetByAggregateIDAndName(stream.Context(),
		tid,
		value.NewPartitionID(request.PartitionId),
		value.NewAggregateID(request.AggregateId),
		value.NewAggregateName(request.Name)) {
		if err != nil {
			return err
		}

		if err := stream.Send(eventToDTO(row)); err != nil {
			return err
		}
	}
	return nil
}

func (esc *EventSourceController) GetAfterGlobalVersion(
	request *pb.GetAfterGlobalVersionRequest,
	stream pb.EventSourceService_GetAfterGlobalVersionServer,
) error {
	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return err
	}

	for row, err := range esc.eventService.GetAfterGlobalVersion(stream.Context(),
		tid,
		value.NewGlobalVersion(request.GlobalVersion),
		dtoToPartitionRange(request.Low, request.High),
		value.Limit(request.Limit)) {
		if err != nil {
			return err
		}

		if err := stream.Send(eventToDTO(row)); err != nil {
			return err
		}
	}
	return nil
}

func (esc *EventSourceController) GetByOutbox(
	request *pb.GetByOutboxRequest,
	stream pb.EventSourceService_GetByOutboxServer,
) error {
	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return err
	}

	oid, err := value.NewOutboxIDFromString(request.OutboxId)
	if err != nil {
		return err
	}

	for row, err := range esc.eventService.GetByOutbox(stream.Context(),
		tid,
		oid,
		value.Limit(request.Limit)) {
		if err != nil {
			return err
		}

		if err := stream.Send(eventToDTO(row)); err != nil {
			return err
		}
	}
	return nil
}

func (esc *EventSourceController) CreateOutbox(
	ctx context.Context,
	request *pb.CreateOutboxRequest,
) (*pb.CreateOutboxResponse, error) {
	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return nil, err
	}

	outboxID, err := esc.eventService.CreateOutbox(ctx,
		tid,
		dtoToPartitionRange(request.Low, request.High))
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

	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return nil, err
	}

	result, err := esc.eventService.GetOutbox(ctx,
		tid,
		oid)
	if err != nil {
		return nil, err
	}

	return outboxToResponse(result), nil
}

func (esc *EventSourceController) UpdateOutboxPosition(
	ctx context.Context,
	request *pb.UpdateOutboxPositionRequest,
) (*pb.UpdateOutboxPositionResponse, error) {
	oid, err := value.NewOutboxIDFromString(request.OutboxId)
	if err != nil {
		return nil, err
	}

	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return nil, err
	}

	if err := esc.eventService.UpdateOutboxPosition(ctx,
		tid,
		oid,
		value.NewGlobalVersion(request.GlobalVersion)); err != nil {
		return nil, err
	}

	return &pb.UpdateOutboxPositionResponse{}, nil
}

func (esc *EventSourceController) CreateTenant(ctx context.Context, request *pb.CreateTenantRequest) (*pb.CreateTenantReponse, error) {
	id, err := esc.tenantService.Create(ctx, service.CreateTenantCommand{
		TenantName: value.NewTenantName(request.Name),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateTenantReponse{
		TenantId: id.String(),
	}, nil
}

func (esc *EventSourceController) DeleteTenant(ctx context.Context, request *pb.DeleteTenantRequest) (*pb.DeleteTenantResponse, error) {
	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return nil, err
	}

	if err := esc.tenantService.Delete(ctx, service.DeleteTenantCommand{
		TenantID: tid,
	}); err != nil {
		return nil, err
	}

	return &pb.DeleteTenantResponse{}, nil
}

func (esc *EventSourceController) RenameTenant(ctx context.Context, request *pb.RenameTenantRequest) (*pb.RenameTenantResponse, error) {
	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return nil, err
	}

	if err := esc.tenantService.Rename(ctx, service.RenameTenantCommand{
		TenantID:   tid,
		TenantName: value.NewTenantName(request.Name),
	}); err != nil {
		return nil, err
	}

	return &pb.RenameTenantResponse{}, nil
}

func (esc *EventSourceController) GetTenant(ctx context.Context, request *pb.GetTenantRequest) (*pb.GetTenantResponse, error) {
	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return nil, err
	}

	t, err := esc.tenantService.Get(ctx, tid)
	if err != nil {
		return nil, err
	}

	return &pb.GetTenantResponse{
		Tenant: tenantToDTO(t),
	}, nil
}

func (esc *EventSourceController) ListOutboxes(request *pb.ListOutboxesRequest, stream pb.EventSourceService_ListOutboxesServer) error {
	tid, err := value.NewTenantIDFromString(request.TenantId)
	if err != nil {
		return err
	}

	for o, err := range esc.eventService.ListOutboxes(stream.Context(), tid) {
		if err != nil {
			return err
		}

		if err := stream.Send(outboxToDTO(o)); err != nil {
			return err
		}
	}
	return nil
}

func (esc *EventSourceController) ListTenants(request *pb.ListTenantsRequest, stream pb.EventSourceService_ListTenantsServer) error {
	for t, err := range esc.tenantService.List(stream.Context(), 0, 0) {
		if err != nil {
			return err
		}

		if err := stream.Send(tenantToDTO(t)); err != nil {
			return err
		}
	}
	return nil
}
