package value

import "github.com/liaradb/liaradb/encoder/raw"

type RequestID struct {
	baseID
}

func NewRequestID() RequestID {
	return RequestID{
		raw.NewBaseID(),
	}
}

func NewRequestIDFromString(value string) (RequestID, error) {
	if id, err := raw.NewBaseIDFromString(value); err != nil {
		return RequestID{}, err
	} else {
		return RequestID{id}, nil
	}
}

const RequestIDIDSize = raw.BaseIDSize
