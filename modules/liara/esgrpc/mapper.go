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
		PartitionID:   liara.PartitionID(dto.PartitionId),
		Schema:        liara.Schema(dto.Schema),
		Metadata: liara.EventMetadata{
			CorrelationID: liara.CorrelationID(dto.Metadata.CorrelationId),
			UserID:        liara.UserID(dto.Metadata.UserId),
			Time:          dto.Metadata.Time.AsTime()},
		Data: dto.Data,
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
		PartitionId:   e.PartitionID.Value(),
		Schema:        e.Schema.String(),
		Metadata: &pb.EventMetadata{
			CorrelationId: e.Metadata.CorrelationID.String(),
			UserId:        e.Metadata.UserID.String(),
			Time:          timestamppb.New(e.Metadata.Time)},
		Data: e.Data,
	}
}

func appendOptionsToDto(options liara.AppendOptions) *pb.AppendOptions {
	return &pb.AppendOptions{
		RequestId:     options.RequestID.String(),
		CorrelationId: options.CorrelationID.String(),
		UserId:        options.UserID.String(),
		Time:          timestamppb.New(options.Time),
	}
}

func appendEventToDto(e liara.AppendEvent) *pb.AppendEvent {
	return &pb.AppendEvent{
		AggregateName: e.AggregateName.String(),
		Id:            e.ID.String(),
		AggregateId:   e.AggregateID.String(),
		Version:       int64(e.Version),
		Name:          e.Name.String(),
		PartitionId:   e.PartitionID.Value(),
		Schema:        e.Schema.String(),
		Data:          e.Data,
	}
}
