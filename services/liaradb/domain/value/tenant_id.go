package value

import (
	"strings"

	"github.com/google/uuid"
)

type TenantID string

func (i TenantID) String() string { return string(i) }

func NewTenantID() TenantID {
	return TenantID(uuid.NewString())
}

func (i TenantID) NewIfEmpty() TenantID {
	id := i.Trim()
	if id == "" {
		return NewTenantID()
	} else {
		return id
	}
}

func (i TenantID) IsEmpty() bool {
	return i.Trim() == ""
}

func (i TenantID) Trim() TenantID {
	return TenantID(strings.TrimSpace(string(i)))
}
