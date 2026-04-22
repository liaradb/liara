package controller

import (
	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func eventToDTO(e *entity.Event) *pb.Event {
	return &pb.Event{
		GlobalVersion: e.GlobalVersion.Value(),
		Id:            e.ID.String(),
		AggregateName: e.AggregateName.String(),
		AggregateId:   e.AggregateID.String(),
		Version:       e.Version.Value(),
		PartitionId:   e.PartitionID.Value(),
		Name:          e.Name.String(),
		Schema:        e.Schema.String(),
		Metadata:      metadataToDto(e.Metadata),
		Data:          e.Data.Value(),
	}
}

func dtoToAppendEvent(dto *pb.AppendEvent) service.AppendEvent {
	return service.AppendEvent{
		ID:            dto.Id,
		AggregateName: value.NewAggregateName(dto.AggregateName),
		AggregateID:   value.NewAggregateID(dto.AggregateId),
		Version:       value.NewVersion(dto.Version),
		Name:          value.NewEventName(dto.Name),
		Schema:        value.NewSchema(dto.Schema),
		Data:          dto.Data,
	}
}

func dtoToAppendOptions(dto *pb.AppendOptions) (service.AppendOptions, error) {
	if dto == nil {
		return service.AppendOptions{}, nil
	}

	var rid *value.RequestID

	if dto.RequestId != "" {
		r, err := value.NewRequestIDFromString(dto.RequestId)
		if err != nil {
			return service.AppendOptions{}, err
		}

		rid = &r
	}

	return service.NewAppendOptions(
		rid,
		value.NewCorrelationID(dto.CorrelationId),
		value.NewClientVersion(dto.ClientVersion),
		value.NewUserID(dto.UserId),
		dto.Time.AsTime(),
	), nil
}

func metadataToDto(m entity.Metadata) *pb.EventMetadata {
	return &pb.EventMetadata{
		CorrelationId: m.CorrelationID().String(),
		UserId:        m.UserID().String(),
		Time:          timestamppb.New(m.Time().Value())}
}

func dtoToPartitionRange(low int32, high int32) value.PartitionRange {
	return value.NewPartitionRange(
		value.NewPartitionID(low),
		value.NewPartitionID(high))
}

func outboxToDTO(o *entity.Outbox) *pb.Outbox {
	low, high := o.PartitionRange().All()

	return &pb.Outbox{
		OutboxId:      o.ID().String(),
		GlobalVersion: o.GlobalVersion().Value(),
		Low:           low.Value(),
		High:          high.Value(),
	}
}

func outboxToResponse(result *entity.Outbox) *pb.GetOutboxResponse {
	low, high := result.PartitionRange().All()

	return &pb.GetOutboxResponse{
		GlobalVersion: result.GlobalVersion().Value(),
		Low:           low.Value(),
		High:          high.Value(),
	}
}

func tenantToDTO(t *entity.Tenant) *pb.Tenant {
	return &pb.Tenant{
		TenantId: t.ID().String(),
		Name:     t.Name().String(),
	}
}

func mapSlice[T any, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, 0, len(slice))
	for _, item := range slice {
		result = append(result, mapper(item))
	}
	return result
}
