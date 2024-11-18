package esmongo

import (
	"encoding/json"
	"errors"
)

type EventDefinition interface {
	ParseEvent(re modelEvent) (Event, error)
}

type EventMapper map[string]EventDefinition

func (em EventMapper) ParseEvent(re modelEvent) (Event, error) {
	parser, ok := em[re.Type]
	if !ok {
		return nil, errors.New("no map")
	}

	return parser.ParseEvent(re)
}

type EventMap[T Event] struct{}

func (EventMap[T]) ParseEvent(re modelEvent) (Event, error) {
	var t T
	return t, json.Unmarshal(re.Data, t)
}
