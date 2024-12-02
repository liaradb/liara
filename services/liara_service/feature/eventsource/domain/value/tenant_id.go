package value

import "github.com/google/uuid"

type TenantID string

func (i TenantID) String() string { return string(i) }

func NewTenantID() TenantID {
	return TenantID(uuid.NewString())
}
