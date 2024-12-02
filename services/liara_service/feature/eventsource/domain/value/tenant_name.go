package value

type TenantName string

func (tn TenantName) String() string { return string(tn) }
