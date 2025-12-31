package controller

import (
	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/domain/value"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func eventToDto(e entity.Event) *pb.Event {
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
		PartitionID:   value.NewPartitionID(uint32(dto.PartitionId)),
		Name:          value.NewEventName(dto.Name),
		Schema:        value.NewSchema(dto.Schema),
		Data:          dto.Data,
	}
}

func dtoToAppendOptions(dto *pb.AppendOptions) service.AppendOptions {
	if dto == nil {
		return service.AppendOptions{}
	}

	return service.AppendOptions{
		RequestID:     value.RequestID(dto.RequestId),
		CorrelationID: value.NewCorrelationID(dto.CorrelationId),
		UserID:        value.NewUserID(dto.UserId),
		Time:          dto.Time.AsTime(),
	}
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
		pids = append(pids, value.NewPartitionID(uint32(p)))
	}
	return value.NewPartitionRange(pids...)
}
