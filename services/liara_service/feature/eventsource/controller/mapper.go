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
		PartitionId:   "",
		Schema:        e.Schema.String(),
		Metadata: &pb.EventMetadata{
			CorrelationId: e.Metadata.CorrelationID.String(),
			UserId:        e.Metadata.UserID.String(),
			Time:          timestamppb.New(e.Metadata.Time),
		},
		Data: e.Data,
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
		PartitionID:   "",
		Metadata: entity.EventMetadata{
			CorrelationID: value.CorrelationID(dto.Metadata.CorrelationId),
			UserID:        value.UserID(dto.Metadata.UserId),
			Time:          dto.Metadata.Time.AsTime()},
		Data: dto.Data,
	}
}
