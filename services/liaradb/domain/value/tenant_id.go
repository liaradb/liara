package value

import (
	"github.com/google/uuid"
	"github.com/liaradb/liaradb/encoder/base"
)

type TenantID struct {
	baseID
}

func NewTenantID() TenantID {
	return TenantID{base.NewID()}
}

func NewTenantIDFromUUID(id uuid.UUID) TenantID {
	return TenantID{
		baseID: base.NewIDFromUUID(id),
	}
}

func NewTenantIDFromString(value string) (TenantID, error) {
	if id, err := base.NewIDFromString(value); err != nil {
		return TenantID{}, err
	} else {
		return TenantID{id}, nil
	}
}

const TenantIDSize = base.BaseIDSize
