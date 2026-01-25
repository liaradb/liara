package value

import (
	"github.com/liaradb/liaradb/encoder/raw"
)

type TenantID struct {
	baseID
}

func NewTenantID() TenantID {
	return TenantID{raw.NewBaseID()}
}

func NewTenantIDFromString(value string) (TenantID, error) {
	if id, err := raw.NewBaseIDFromString(value); err != nil {
		return TenantID{}, err
	} else {
		return TenantID{id}, nil
	}
}

const TenantIDSize = raw.BaseIDSize
