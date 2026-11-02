package tablename

import "github.com/liaradb/liaradb/domain/value"

type TenantFolder struct {
	tenantID value.TenantID
}

func NewTenantFolder(tid value.TenantID) TenantFolder {
	return TenantFolder{tid}
}

func NewTenantFolderFromString(s string) (TenantFolder, error) {
	tid, err := value.NewTenantIDFromString(s)
	if err != nil {
		return TenantFolder{}, err
	}

	return TenantFolder{tid}, nil
}

func (tf TenantFolder) TenantID() value.TenantID { return tf.tenantID }
func (tf TenantFolder) String() string           { return tf.tenantID.String() }
func (tf TenantFolder) Table(name string) Table  { return newTable(tf, name) }
