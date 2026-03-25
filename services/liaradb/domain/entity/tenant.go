package entity

import "github.com/liaradb/liaradb/domain/value"

const (
	TenantSize = value.TenantIDSize +
		value.VersionSize +
		value.TenantNameSize
)

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
	t.version.Increment()
	return nil
}

func (t *Tenant) Write(data []byte) ([]byte, bool) {
	data0 := t.id.WriteData(data)
	data1, ok := t.version.WriteData(data0)
	if !ok {
		return nil, false
	}

	return t.name.WriteData(data1), true
}

func (t *Tenant) Read(data []byte) ([]byte, bool) {
	data0 := t.id.ReadData(data)
	data1, ok := t.version.ReadData(data0)
	if !ok {
		return nil, false
	}

	return t.name.ReadData(data1), true
}
