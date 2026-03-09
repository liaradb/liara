package value

import "github.com/liaradb/liaradb/encoder/base"

type TenantID struct {
	baseID
}

func NewTenantID() TenantID {
	return TenantID{base.NewID()}
}

func NewTenantIDFromString(value string) (TenantID, error) {
	if id, err := base.NewIDFromString(value); err != nil {
		return TenantID{}, err
	} else {
		return TenantID{id}, nil
	}
}

const TenantIDSize = base.BaseIDSize
