package value

import "github.com/liaradb/liaradb/encoder/base"

type OutboxID struct {
	baseID
}

func NewOutboxID() OutboxID {
	return OutboxID{base.NewID()}
}

func NewOutboxIDFromString(value string) (OutboxID, error) {
	if id, err := base.NewIDFromString(value); err != nil {
		return OutboxID{}, err
	} else {
		return OutboxID{id}, nil
	}
}

const OutboxIDSize = base.BaseIDSize
