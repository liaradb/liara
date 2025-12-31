package value

import "github.com/liaradb/liaradb/encoder/raw"

type RowID struct {
	baseID
}

func NewRowID() RowID {
	return RowID{raw.NewBaseID()}
}

func NewRowIDFromString(value string) (RowID, error) {
	if id, err := raw.NewBaseIDFromString(value); err != nil {
		return RowID{}, err
	} else {
		return RowID{id}, nil
	}
}

const RowIDSize = raw.BaseIDSize
