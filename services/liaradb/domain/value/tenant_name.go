package value

import "github.com/liaradb/liaradb/encoder/base"

type TenantName struct {
	baseShortString
}

func NewTenantName(value string) TenantName {
	return TenantName{base.NewShortString(value)}
}
