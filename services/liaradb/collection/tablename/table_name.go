package tablename

import (
	"github.com/liaradb/liaradb/encoder/base"
	"github.com/liaradb/liaradb/storage/link"
)

type TableName interface {
	Collection() link.FileName
	Index(i base.Uint32) link.FileName
	String() string
}
