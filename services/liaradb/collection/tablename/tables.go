package tablename

import "github.com/liaradb/liaradb/domain/value"

const (
	managerTable = "tables"
	tenantTable  = "tenants"
)

var (
	Tenant  = NewTenantFolder(value.DefaultTenantID).GlobalTable(tenantTable)
	Manager = NewTenantFolder(value.DefaultTenantID).GlobalTable(managerTable)
)
