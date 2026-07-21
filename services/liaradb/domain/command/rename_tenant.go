package command

import "github.com/liaradb/liaradb/domain/value"

type RenameTenant struct {
	TenantID   value.TenantID
	TenantName value.TenantName
}
