package value

import "github.com/liaradb/liaradb/encoder/base"

type RequestID struct {
	baseID
}

func NewRequestID() RequestID {
	return RequestID{
		base.NewBaseID(),
	}
}

func NewRequestIDFromString(value string) (RequestID, error) {
	if id, err := base.NewBaseIDFromString(value); err != nil {
		return RequestID{}, err
	} else {
		return RequestID{id}, nil
	}
}

const RequestIDIDSize = base.BaseIDSize
