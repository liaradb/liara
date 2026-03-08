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
	time value.Time
}

func NewRequestLog(
	id value.RequestID,
	t value.Time,
) *RequestLog {
	return &RequestLog{
		id:   id,
		time: t,
	}
}

func RestoreRequestLog(
	id value.RequestID,
	t value.Time,
) *RequestLog {
	return &RequestLog{
		id:   id,
		time: t,
	}
}

func (rl *RequestLog) ID() value.RequestID { return rl.id }
func (rl *RequestLog) Time() value.Time    { return rl.time }

func (rl *RequestLog) Compare(b *RequestLog) bool {
	if rl == b {
		return true
	}

	return rl.id == b.id && rl.time.Compare(b.time.Value()) == 0
}

func (rl *RequestLog) Write(data []byte) []byte {
	data0 := rl.id.WriteData(data)
	t := rl.time.UnixMicro()
	return scan.SetInt64(data0, t)
}

func (rl *RequestLog) Read(data []byte) []byte {
	data0 := rl.id.ReadData(data)
	t, data1 := scan.Int64(data0)
	rl.time = value.NewTime(time.UnixMicro(t).UTC())
	return data1
}
