package fixed

import (
	"context"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
)

type Transaction interface {
	span.Log

	// TODO: Don't use TableName
	// TODO: Maybe use an interface for string
	XLock(context.Context, tablename.TableName, string) error
	SLock(context.Context, tablename.TableName, string) error
}
