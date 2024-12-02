package entity

import "github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"

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
	return &Tenant{
		id:      id,
		version: 0,
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
