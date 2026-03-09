package schema

import "github.com/liaradb/liaradb/encoder/base"

type baseUUID = base.BaseID

type ID struct {
	baseUUID
}

func NewID() ID {
	return ID{base.NewBaseID()}
}
