package value

import "strings"

type TenantName string

func (tn TenantName) String() string { return string(tn) }

func (tn TenantName) IsEmpty() bool {
	return tn.trim() == ""
}

func (tn TenantName) trim() TenantName {
	return TenantName(strings.TrimSpace(string(tn)))
}
