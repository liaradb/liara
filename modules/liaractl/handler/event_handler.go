package handler

import (
	"context"
	"encoding/json"
	"fmt"

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
	count := 0
	for event, err := range eh.es.GetAfterGlobalVersion(context.Background(), 0, nil, 0) {
		if err != nil {
			return err
		}

		var data = make(map[string]any)
		_ = json.Unmarshal(event.Data, &data)
		result, _ := json.MarshalIndent(data, "", "    ")

		fmt.Printf("%v,\n", string(result))
		count++
	}
	if count == 0 {
		fmt.Println("no events")
	}
	return nil
}
