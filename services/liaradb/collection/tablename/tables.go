package tablename

const (
	managerTable = "tables"
	tenantTable  = "tenants"
)

var (
	Tenant  = NewFromString(tenantTable)
	Manager = NewFromString(managerTable)
)
