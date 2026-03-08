package entity

import (
	"testing"
	"time"

	"github.com/liaradb/liaradb/domain/value"
)

func TestRequestLog(t *testing.T) {
	rlid := value.NewRequestID()
	tm := time.Now().UTC()
	rl := NewRequestLog(rlid, tm)

	if i := rl.ID(); i != rlid {
		t.Errorf("incorrect id: %v, expected: %v", i, rlid)
	}

	if tm2 := rl.Time().Value(); tm2 != tm {
		t.Errorf("incorrect time: %v, expected: %v", tm2, tm)
	}
}

func TestRequestLog_RestoreRequestLog(t *testing.T) {
	rlid := value.NewRequestID()
	tm := time.Now().UTC()
	rl := RestoreRequestLog(rlid, tm)

	if i := rl.ID(); i != rlid {
		t.Errorf("incorrect id: %v, expected: %v", i, rlid)
	}

	if tm2 := rl.Time().Value(); tm2 != tm {
		t.Errorf("incorrect time: %v, expected: %v", tm2, tm)
	}
}

func TestRequestLog_ReadWrite(t *testing.T) {
	rl := NewRequestLog(
		value.NewRequestID(),
		time.Now())

	data := make([]byte, RequestLogSize+2)
	data0 := rl.Write(data)

	if l := len(data0); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	rl1 := &RequestLog{}
	data1 := rl1.Read(data)
	if l := len(data1); l != 2 {
		t.Errorf("incorrect length: %v, expected: %v", l, 2)
	}

	if *rl1 != *rl {
		t.Errorf("incorrect result: %v, expected: %v", *rl1, *rl)
	}
}
