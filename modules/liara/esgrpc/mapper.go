package esgrpc

import (
	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func DtoToEvent(dto *pb.Event) liara.Event {
	return liara.Event{
		GlobalVersion: liara.GlobalVersion(dto.GlobalVersion),
		AggregateName: liara.AggregateName(dto.AggregateName),
		ID:            liara.EventID(dto.Id),
		AggregateID:   liara.AggregateID(dto.AggregateId),
		Version:       liara.Version(dto.Version),
		Name:          liara.EventName(dto.Name),
		CorrelationID: liara.CorrelationID(dto.CorrelationId),
		UserID:        liara.UserID(dto.UserId),
		Time:          dto.Time.AsTime(),
		Schema:        liara.Schema(dto.Schema),
		Data:          dto.Data,
	}
}

func EventToDto(e liara.Event) *pb.Event {
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

func DtoToAppendEvent(dto *pb.AppendEvent) liara.AppendEvent {
	return liara.AppendEvent{
		AggregateName: liara.AggregateName(dto.AggregateName),
		ID:            liara.EventID(dto.Id),
		AggregateID:   liara.AggregateID(dto.AggregateId),
		Version:       liara.Version(dto.Version),
		Name:          liara.EventName(dto.Name),
		CorrelationID: liara.CorrelationID(dto.CorrelationId),
		UserID:        liara.UserID(dto.UserId),
		Time:          dto.Time.AsTime(),
		Schema:        liara.Schema(dto.Schema),
		Data:          dto.Data,
	}
}

func AppendEventToDto(e liara.AppendEvent) *pb.AppendEvent {
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
