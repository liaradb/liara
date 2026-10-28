package transaction

import (
	"fmt"

	"github.com/liaradb/liaradb/collection/tablename"
)

type ItemID struct {
	tn    tablename.TableName
	value string
}

func newItemID(tn tablename.TableName, value string) ItemID {
	return ItemID{
		tn:    tn,
		value: value,
	}
}

func (i ItemID) String() string {
	return fmt.Sprintf("%v|%v", i.tn.String(), i.value)
}
