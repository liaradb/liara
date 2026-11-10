package controller

import (
	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func eventToDTO(e *entity.Event) *pb.Event {
	return &pb.Event{
		GlobalVersion: e.GlobalVersion.Value(),
		SubVersion:    e.SubVersion.Value(),
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

func dtoToAppendEvent(dto *pb.AppendEvent) command.EventOptions {
	return command.EventOptions{
		ID:            dto.Id,
		AggregateName: value.NewAggregateName(dto.AggregateName),
		AggregateID:   value.NewAggregateID(dto.AggregateId),
		Version:       value.NewVersion(dto.Version),
		PartitionID:   value.NewPartitionID(dto.PartitionId),
		Name:          value.NewEventName(dto.Name),
		Schema:        value.NewSchema(dto.Schema),
		Data:          dto.Data,
	}
}

func dtoToAppendOptions(dto *pb.AppendOptions) (command.AppendOptions, error) {
	if dto == nil {
		return command.AppendOptions{}, nil
	}

	var rid *value.RequestID

	if dto.RequestId != "" {
		r, err := value.NewRequestIDFromString(dto.RequestId)
		if err != nil {
			return command.AppendOptions{}, err
		}

		rid = &r
	}

	return command.NewAppendOptions(
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

func outboxToDTO(o *entity.Outbox) *pb.Outbox {
	return &pb.Outbox{
		OutboxId:      o.ID().String(),
		GlobalVersion: o.GlobalVersion().Value(),
		PartitionId:   o.PartitionID().Value(),
	}
}

func outboxToResponse(result *entity.Outbox) *pb.GetOutboxResponse {
	return &pb.GetOutboxResponse{
		GlobalVersion: result.GlobalVersion().Value(),
		PartitionId:   result.PartitionID().Value(),
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
