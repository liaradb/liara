package liara

type Tenant struct {
	ID   TenantID
	Name TenantName
}

type TenantName string

func (t TenantName) String() string { return string(t) }
