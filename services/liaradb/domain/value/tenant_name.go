package value

import "github.com/liaradb/liaradb/encoder/raw"

const TenantNameSize = 256

type TenantName struct {
	baseString
}

func NewTenantName(value string) TenantName {
	return TenantName{raw.BaseString(value)}
}

func (tn TenantName) WriteData(data []byte) []byte {
	return tn.baseString.WriteData(data, TenantNameSize)
}

func (tn *TenantName) ReadData(data []byte) []byte {
	return tn.baseString.ReadData(data, TenantNameSize)
}
