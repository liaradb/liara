package value

import "github.com/liaradb/liaradb/encoder/base"

type TenantName struct {
	baseShortString[base.ShortString255]
}

func NewTenantName(value string) TenantName {
	return TenantName{base.NewShortString[base.ShortString255](value)}
}
