package entity

import (
	"time"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type RequestLog struct {
	ID   value.RequestID
	Time time.Time
}
