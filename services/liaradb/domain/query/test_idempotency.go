package query

import "github.com/liaradb/liaradb/domain/value"

type TestIdempotency struct {
	TenantID  value.TenantID
	RequestID value.RequestID
}
