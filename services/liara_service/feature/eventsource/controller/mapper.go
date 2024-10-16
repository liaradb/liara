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
		CorrelationId: e.CorrelationID.String(),
		UserId:        e.UserID.String(),
		Time:          timestamppb.New(e.Time),
		Schema:        e.Schema.String(),
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
		CorrelationID: value.CorrelationID(dto.CorrelationId),
		UserID:        value.UserID(dto.UserId),
		Time:          dto.Time.AsTime(),
		Schema:        value.Schema(dto.Schema),
		Data:          dto.Data,
	}
}
