package value

import "github.com/liaradb/liaradb/encoder/base"

type RowID struct {
	baseID
}

func NewRowID() RowID {
	return RowID{base.NewID()}
}

func NewRowIDFromString(value string) (RowID, error) {
	if id, err := base.NewIDFromString(value); err != nil {
		return RowID{}, err
	} else {
		return RowID{id}, nil
	}
}

const RowIDSize = base.BaseIDSize
