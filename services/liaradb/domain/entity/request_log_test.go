package entity

import (
	"testing"
	"time"

	"github.com/liaradb/liaradb/domain/value"
)

func TestRequestLog(t *testing.T) {
	rlid := value.NewRequestID()
	tm := value.NewTime(time.Now())
	rl := NewRequestLog(rlid, tm)

	if i := rl.ID(); i != rlid {
		t.Errorf("incorrect id: %v, expected: %v", i, rlid)
	}

	if tm2 := rl.Time(); tm2 != tm {
		t.Errorf("incorrect time: %v, expected: %v", tm2, tm)
	}
}

func TestRequestLog_RestoreRequestLog(t *testing.T) {
	rlid := value.NewRequestID()
	tm := value.NewTime(time.Now())
	rl := RestoreRequestLog(rlid, tm)

	if i := rl.ID(); i != rlid {
		t.Errorf("incorrect id: %v, expected: %v", i, rlid)
	}

	if tm2 := rl.Time(); tm2 != tm {
		t.Errorf("incorrect time: %v, expected: %v", tm2, tm)
	}
}

func TestRequestLog_ReadWrite(t *testing.T) {
	rl := NewRequestLog(
		value.NewRequestID(),
		value.NewTime(time.Now()))

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

	if !rl1.Compare(rl) {
		t.Errorf("incorrect result: %v, expected: %v", *rl1, *rl)
	}
}

func TestRequestLog_Compare(t *testing.T) {
	pointer := &RequestLog{}
	rid := value.NewRequestID()
	tm := value.NewTime(time.Now())
	for message, c := range map[string]struct {
		skip  bool
		a     *RequestLog
		b     *RequestLog
		equal bool
	}{
		"should equal zero": {
			a:     &RequestLog{},
			b:     &RequestLog{},
			equal: true,
		},
		"should equal pointer": {
			a:     pointer,
			b:     pointer,
			equal: true,
		},
		"should equal same values": {
			a:     NewRequestLog(rid, tm),
			b:     NewRequestLog(rid, tm),
			equal: true,
		},
		"should not equal different values": {
			a:     NewRequestLog(rid, tm),
			b:     NewRequestLog(value.NewRequestID(), value.NewTime(time.Now().Add(1*time.Second))),
			equal: false,
		}} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			if c.a.Compare(c.b) != c.equal {
				if c.equal {
					t.Error("should equal")
				} else {
					t.Error("should not equal")
				}
			}
		})
	}
}
