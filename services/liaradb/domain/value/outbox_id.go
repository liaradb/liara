package value

import "github.com/liaradb/liaradb/encoder/raw"

type OutboxID struct {
	baseID
}

func NewOutboxID() OutboxID {
	return OutboxID{raw.NewBaseID()}
}

func NewOutboxIDFromString(value string) (OutboxID, error) {
	if id, err := raw.NewBaseIDFromString(value); err != nil {
		return OutboxID{}, err
	} else {
		return OutboxID{id}, nil
	}
}

const OutboxIDSize = raw.BaseIDSize
