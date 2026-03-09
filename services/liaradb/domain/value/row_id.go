package value

import "github.com/liaradb/liaradb/encoder/base"

type RowID struct {
	baseID
}

func NewRowID() RowID {
	return RowID{base.NewBaseID()}
}

func NewRowIDFromString(value string) (RowID, error) {
	if id, err := base.NewBaseIDFromString(value); err != nil {
		return RowID{}, err
	} else {
		return RowID{id}, nil
	}
}

const RowIDSize = base.BaseIDSize
