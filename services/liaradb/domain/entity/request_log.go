package entity

import (
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/encoder/scan"
)

const (
	RequestLogSize = value.RequestIDIDSize +
		raw.TimeSize
)

type RequestLog struct {
	id   value.RequestID
	time time.Time // TODO: Change to raw.BaseTime
}

func NewRequestLog(
	id value.RequestID,
	t time.Time,
) *RequestLog {
	return &RequestLog{
		id:   id,
		time: t.Truncate(time.Microsecond).UTC(),
	}
}

func RestoreRequestLog(
	id value.RequestID,
	t time.Time,
) *RequestLog {
	return &RequestLog{
		id:   id,
		time: t,
	}
}

func (rl *RequestLog) ID() value.RequestID { return rl.id }
func (rl *RequestLog) Time() time.Time     { return rl.time }

func (rl *RequestLog) Write(data []byte) []byte {
	data0 := rl.id.WriteData(data)
	t := rl.time.UnixMicro()
	return scan.SetInt64(data0, t)
}

func (rl *RequestLog) Read(data []byte) []byte {
	data0 := rl.id.ReadData(data)
	t, data1 := scan.Int64(data0)
	rl.time = time.UnixMicro(t).UTC()
	return data1
}
