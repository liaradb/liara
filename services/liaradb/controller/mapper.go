package controller

import (
	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func eventToDto(e *entity.Event) *pb.Event {
	return &pb.Event{
		GlobalVersion: int64(e.GlobalVersion.Value()),
		Id:            e.ID.String(),
		AggregateName: e.AggregateName.String(),
		AggregateId:   e.AggregateID.String(),
		Version:       int64(e.Version.Value()),
		PartitionId:   int32(e.PartitionID.Value()),
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
		Version:       value.NewVersion(uint64(dto.Version)),
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
		value.NewUserID(dto.UserId),
		dto.Time.AsTime(),
	), nil
}

func metadataToDto(m entity.Metadata) *pb.EventMetadata {
	return &pb.EventMetadata{
		CorrelationId: m.CorrelationID.String(),
		UserId:        m.UserID.String(),
		Time:          timestamppb.New(m.Time.Value())}
}

func dtoToPartitionRange(partitionIDs []int32) value.PartitionRange {
	pids := make([]value.PartitionID, 0, len(partitionIDs))
	for _, p := range partitionIDs {
		pids = append(pids, value.NewPartitionID(p))
	}
	return value.NewPartitionRange(pids...)
}
