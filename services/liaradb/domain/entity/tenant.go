package entity

import "github.com/liaradb/liaradb/domain/value"

type Tenant struct {
	id      value.TenantID
	version value.Version
	name    value.TenantName
}

func (t *Tenant) ID() value.TenantID     { return t.id }
func (t *Tenant) Version() value.Version { return t.version }
func (t *Tenant) Name() value.TenantName { return t.name }

func NewTenant(
	id value.TenantID,
	name value.TenantName,
) *Tenant {
	if !id.IsEmpty() {
		if name.IsEmpty() {
			name = value.TenantName(id.Trim())
		}
	} else {
		id = value.NewTenantID()
	}

	return &Tenant{
		id:      id,
		version: value.NewVersion(0),
		name:    name,
	}
}

func RestoreTenant(
	id value.TenantID,
	version value.Version,
	name value.TenantName,
) *Tenant {
	return &Tenant{
		id:      id,
		version: version,
		name:    name,
	}
}

func (t *Tenant) Rename(name value.TenantName) error {
	t.name = name
	return nil
}
