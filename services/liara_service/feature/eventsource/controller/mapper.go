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
		AggregateName: e.AggregateName.String(),
		Id:            e.ID.String(),
		AggregateId:   e.AggregateID.String(),
		Version:       int64(e.Version),
		Name:          e.Name.String(),
		PartitionId:   e.PartitionID.String(),
		Schema:        e.Schema.String(),
		Metadata:      metadataToDto(e.Metadata),
		Data:          e.Data,
	}
}

func dtoToAppendEvent(dto *pb.AppendEvent) service.AppendEvent {
	return service.AppendEvent{
		AggregateName: value.AggregateName(dto.AggregateName),
		ID:            value.EventID(dto.Id),
		AggregateID:   value.AggregateID(dto.AggregateId),
		Version:       value.Version(dto.Version),
		Name:          value.EventName(dto.Name),
		Schema:        value.Schema(dto.Schema),
		PartitionID:   value.PartitionID(dto.PartitionId),
		Metadata:      dtoToMetadata(dto.Metadata),
		Data:          dto.Data,
	}
}

func metadataToDto(m entity.EventMetadata) *pb.EventMetadata {
	return &pb.EventMetadata{
		CorrelationId: m.CorrelationID.String(),
		UserId:        m.UserID.String(),
		Time:          timestamppb.New(m.Time)}
}

func dtoToMetadata(dto *pb.EventMetadata) entity.EventMetadata {
	if dto == nil {
		return entity.EventMetadata{}
	}

	return entity.EventMetadata{
		CorrelationID: value.CorrelationID(dto.CorrelationId),
		UserID:        value.UserID(dto.UserId),
		Time:          dto.Time.AsTime()}
}
