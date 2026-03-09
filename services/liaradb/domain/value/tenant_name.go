package value

import "github.com/liaradb/liaradb/encoder/base"

const TenantNameSize = 256

type TenantName struct {
	baseString
}

func NewTenantName(value string) TenantName {
	return TenantName{base.String(value)}
}

func (tn TenantName) WriteData(data []byte) []byte {
	return tn.baseString.WriteData(data, TenantNameSize)
}

func (tn *TenantName) ReadData(data []byte) []byte {
	return tn.baseString.ReadData(data, TenantNameSize)
}
