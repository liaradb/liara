package entity

import (
	"io"
	"testing"
	"time"

	"github.com/liaradb/liaradb/domain/value"
)

func TestRow(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var row = Row{}
	row.SetData(value.NewData([]byte{}))

	if err := row.Write(w); err != nil {
		t.Fatal(err)
	}

	// size := w.Len()
	// if s := lsn.Size(); s != size {
	// 	t.Errorf("incorrect size: %v, expected: %v", s, size)
	// }

	var row2 Row
	if err := row2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	// Data comparison doesn't allow nil slice
	if !row.Compare(&row2) {
		t.Errorf("incorrect value: %v, expected: %v", row2, row)
	}
}

func TestRow_Compare(t *testing.T) {
	t.Parallel()

	rid := value.NewRowID()
	v := value.NewVersion(1)
	pid := value.NewPartitionID(2)
	rn := value.NewRowName("rowname")
	sch := value.NewSchema("schema")
	uid := value.NewUserID("userid")
	cid := value.NewCorrelationID("correlationid")
	cver := "clientversion"
	tm := value.NewTime(time.Now())

	for message, c := range map[string]struct {
		skip  bool
		a     Row
		b     Row
		equal bool
	}{
		"should equal zero": {
			a:     Row{},
			b:     Row{},
			equal: true,
		},
		"should equal same values": {
			a: *NewRow(
				rid,
				v,
				pid,
				rn,
				sch,
				Metadata{
					UserID:        uid,
					CorrelationID: cid,
					ClientVersion: cver,
					Time:          tm,
				},
				value.NewData([]byte{1, 2, 3}),
			),
			b: *NewRow(
				rid,
				v,
				pid,
				rn,
				sch,
				Metadata{
					UserID:        uid,
					CorrelationID: cid,
					ClientVersion: cver,
					Time:          tm,
				},
				value.NewData([]byte{1, 2, 3}),
			),
			equal: true,
		},
		"should not equal different values": {
			a: *NewRow(
				rid,
				v,
				pid,
				rn,
				sch,
				Metadata{
					UserID:        uid,
					CorrelationID: cid,
					ClientVersion: cver,
					Time:          tm,
				},
				value.NewData([]byte{1, 2, 3}),
			),
			b: *NewRow(
				value.NewRowID(),
				v,
				pid,
				rn,
				sch,
				Metadata{
					UserID:        uid,
					CorrelationID: cid,
					ClientVersion: cver,
					Time:          tm,
				},
				value.NewData([]byte{1, 2, 3}),
			),
			equal: false,
		},
		"should not equal different data": {
			a: *NewRow(
				rid,
				v,
				pid,
				rn,
				sch,
				Metadata{
					UserID:        uid,
					CorrelationID: cid,
					ClientVersion: cver,
					Time:          tm,
				},
				value.NewData([]byte{1, 2, 3}),
			),
			b: *NewRow(
				rid,
				v,
				pid,
				rn,
				sch,
				Metadata{
					UserID:        uid,
					CorrelationID: cid,
					ClientVersion: cver,
					Time:          tm,
				},
				value.NewData([]byte{3, 2, 1}),
			),
			equal: false,
		},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if c.skip {
				t.Skip()
			}

			if c.a.Compare(&c.b) != c.equal {
				if c.equal {
					t.Error("should equal")
				} else {
					t.Error("should not equal")
				}
			}
		})
	}
}
