package table

import "github.com/liaradb/liaradb/encoder/raw"

type baseUUID = raw.BaseID

type ID struct {
	baseUUID
}
