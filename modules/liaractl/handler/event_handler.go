package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cardboardrobots/liara"
	"github.com/cardboardrobots/liara/esgrpc"
)

type eventHandler struct {
	es *esgrpc.EventSourceGRPC
}

func (eh *eventHandler) handle(cmd command) error {
	switch cmd {
	case commandEventList:
		return eh.listEvents()
	default:
		return errNoCommand
	}
}

func (eh *eventHandler) listEvents() error {
	w := json.NewEncoder(os.Stdout)
	w.SetIndent("", "  ")
	count := 0
	for event, err := range eh.es.GetAfterGlobalVersion(context.Background(), 0, nil, 0) {
		if err != nil {
			return err
		}

		w.Encode(eventToRecord(event))
		count++
	}
	if count == 0 {
		fmt.Println("no events")
	}
	return nil
}

func eventToRecord(event liara.Event) Event {
	var data = make(map[string]any)
	_ = json.Unmarshal(event.Data, &data)

	return Event{
		GlobalVersion: event.GlobalVersion,
		ID:            event.ID,
		AggregateName: event.AggregateName,
		AggregateID:   event.AggregateID,
		Version:       event.Version,
		PartitionID:   event.PartitionID,
		Name:          event.Name,
		Schema:        event.Schema,
		Metadata:      EventMetadata(event.Metadata),
		Data:          data,
	}
}

type Event struct {
	GlobalVersion liara.GlobalVersion `json:"globalVersion"`
	ID            liara.EventID       `json:"id"`
	AggregateName liara.AggregateName `json:"aggregateName"`
	AggregateID   liara.AggregateID   `json:"aggregateId"`
	Version       liara.Version       `json:"version"`
	PartitionID   liara.PartitionID   `json:"partitionId"`
	Name          liara.EventName     `json:"name"`
	Schema        liara.Schema        `json:"schema"`
	Metadata      EventMetadata       `json:"metadata"`
	Data          any                 `json:"data"`
}

type EventMetadata struct {
	UserID        liara.UserID        `json:"userId"`
	CorrelationID liara.CorrelationID `json:"correlationId"`
	Time          time.Time           `json:"time"`
}
