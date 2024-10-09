package esgrpc

import (
	"github.com/cardboardrobots/eventsource/entity"
	"github.com/cardboardrobots/eventsource/value"
	pb "github.com/cardboardrobots/eventsource_go/generated"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func DtoToEvent(dto *pb.Event) entity.Event {
	return entity.Event{
		GlobalVersion: value.GlobalVersion(dto.GlobalVersion),
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

func EventToDto(e entity.Event) *pb.Event {
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

func DtoToAppendEvent(dto *pb.AppendEvent) entity.AppendEvent {
	return entity.AppendEvent{
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

func AppendEventToDto(e entity.AppendEvent) *pb.AppendEvent {
	return &pb.AppendEvent{
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
