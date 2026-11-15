package fixed

import (
	"context"

	"github.com/liaradb/liaradb/collection/log"
	"github.com/liaradb/liaradb/collection/tablename"
)

type Transaction interface {
	log.Log

	XLock(context.Context, tablename.TableName, string) error
	SLock(context.Context, tablename.TableName, string) error
}
