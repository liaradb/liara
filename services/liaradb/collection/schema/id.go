package schema

import "github.com/liaradb/liaradb/encoder/raw"

type baseUUID = raw.BaseID

type ID struct {
	baseUUID
}

func NewID() ID {
	return ID{raw.NewBaseID()}
}
