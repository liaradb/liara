package schema

import "github.com/liaradb/liaradb/encoder/base"

type baseUUID = base.ID

type ID struct {
	baseUUID
}

func NewID() ID {
	return ID{base.NewID()}
}
