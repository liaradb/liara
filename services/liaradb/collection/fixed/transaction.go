package fixed

import (
	"context"

	"github.com/liaradb/liaradb/collection/span"
	"github.com/liaradb/liaradb/collection/tablename"
)

type Transaction interface {
	span.Log

	XLock(context.Context, tablename.TableName, string) error
	SLock(context.Context, tablename.TableName, string) error
}
