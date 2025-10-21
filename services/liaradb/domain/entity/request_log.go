package entity

import (
	"time"

	"github.com/liaradb/liaradb/domain/value"
)

type RequestLog struct {
	ID   value.RequestID
	Time time.Time
}
