package controller

import (
	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/service"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func eventToDto(e entity.Event) *pb.Event {
	return &pb.Event{
		GlobalVersion: int64(e.GlobalVersion),
		Id:            e.ID.String(),
		AggregateName: e.AggregateName.String(),
		AggregateId:   e.AggregateID.String(),
		Version:       int64(e.Version),
		PartitionId:   e.PartitionID.Value(),
		Name:          e.Name.String(),
		Schema:        e.Schema.String(),
		Metadata:      metadataToDto(e.Metadata),
		Data:          e.Data,
	}
}

func dtoToAppendEvent(dto *pb.AppendEvent) service.AppendEvent {
	return service.AppendEvent{
		ID:            value.EventID(dto.Id),
		AggregateName: value.AggregateName(dto.AggregateName),
		AggregateID:   value.AggregateID(dto.AggregateId),
		Version:       value.Version(dto.Version),
		PartitionID:   value.PartitionID(dto.PartitionId),
		Name:          value.EventName(dto.Name),
		Schema:        value.Schema(dto.Schema),
		Data:          dto.Data,
	}
}

func dtoToAppendOptions(dto *pb.AppendOptions) service.AppendOptions {
	if dto == nil {
		return service.AppendOptions{}
	}

	return service.AppendOptions{
		RequestID:     value.RequestID(dto.RequestId),
		CorrelationID: value.CorrelationID(dto.CorrelationId),
		UserID:        value.UserID(dto.UserId),
		Time:          dto.Time.AsTime(),
	}
}

func metadataToDto(m entity.EventMetadata) *pb.EventMetadata {
	return &pb.EventMetadata{
		CorrelationId: m.CorrelationID.String(),
		UserId:        m.UserID.String(),
		Time:          timestamppb.New(m.Time)}
}

func dtoToPartitionRange(partitionIDs []int32) value.PartitionRange {
	pids := make([]value.PartitionID, 0, len(partitionIDs))
	for _, p := range partitionIDs {
		pids = append(pids, value.PartitionID(p))
	}
	return value.NewPartitionRange(pids...)
}
