package value

import "github.com/liaradb/liaradb/encoder/base"

type TenantID struct {
	baseID
}

func NewTenantID() TenantID {
	return TenantID{base.NewBaseID()}
}

func NewTenantIDFromString(value string) (TenantID, error) {
	if id, err := base.NewBaseIDFromString(value); err != nil {
		return TenantID{}, err
	} else {
		return TenantID{id}, nil
	}
}

const TenantIDSize = base.BaseIDSize
