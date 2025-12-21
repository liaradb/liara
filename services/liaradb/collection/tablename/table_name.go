package tablename

import "fmt"

type TableName struct {
	n string
}

func New(n string) TableName {
	return TableName{n}
}

func (tn *TableName) KeyValue() string {
	return fmt.Sprintf("%v.kv", tn.n)
}

func (tn *TableName) EventLog() string {
	return fmt.Sprintf("%v.el", tn.n)
}

func (tn *TableName) Index(i int) string {
	return fmt.Sprintf("%v--%v.idx", tn.n, i)
}
