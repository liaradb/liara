package entity

import (
	"io"
	"slices"
	"testing"
	"time"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
)

func TestRow_Default(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	var row = Row{}
	row.SetData(value.NewData([]byte{}))

	if err := row.Write(w); err != nil {
		t.Fatal(err)
	}

	size := w.Len() - raw.HeaderSize
	if s := row.Size(); s != size {
		t.Errorf("incorrect size: %v, expected: %v", s, size)
	}

	var row2 Row
	if err := row2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	// Data comparison doesn't allow nil slice
	if !row.Compare(&row2) {
		t.Errorf("incorrect value: %v, expected: %v", row2, row)
	}
}

func TestRow_NewRow(t *testing.T) {
	t.Parallel()

	r, w := newReaderWriter()

	row := NewRow(
		value.NewRowID(),
		value.NewVersion(1),
		value.NewPartitionID(2),
		value.NewRowName("name"),
		value.NewSchema("schema"),
		NewMetadata(
			value.NewUserID("userId"),
			value.NewCorrelationID("correlationID"),
			value.NewClientVersion("clientVersion"),
			value.NewTime(time.Time{}),
		),
		value.NewData([]byte{1, 2, 3, 4}),
	)

	if err := row.Write(w); err != nil {
		t.Fatal(err)
	}

	var row2 Row
	if err := row2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	// Data comparison doesn't allow nil slice
	if !row.Compare(&row2) {
		t.Errorf("incorrect value: %v, expected: %v", row2, row)
	}
}

func TestRow__Getters(t *testing.T) {
	t.Parallel()

	rid := value.NewRowID()
	version := value.NewVersion(1)
	pid := value.NewPartitionID(2)
	name := value.NewRowName("name")
	schema := value.NewSchema("schema")
	userID := value.NewUserID("userId")
	correlationID := value.NewCorrelationID("correlationID")
	clientVersion := value.NewClientVersion("clientVersion")
	tm := value.NewTime(time.Time{})
	data := value.NewData([]byte{1, 2, 3, 4})

	row := NewRow(
		rid,
		version,
		pid,
		name,
		schema,
		NewMetadata(
			userID,
			correlationID,
			clientVersion,
			tm,
		),
		data,
	)

	if v := row.ID(); v != rid {
		t.Errorf("incorrect id: %v, expected: %v", v, rid)
	}

	if v := row.Version(); v != version {
		t.Errorf("incorrect version: %v, expected: %v", v, version)
	}

	if v := row.PartitionID(); v != pid {
		t.Errorf("incorrect partition id: %v, expected: %v", v, pid)
	}

	if v := row.Name(); v != name {
		t.Errorf("incorrect name: %v, expected: %v", v, name)
	}

	if v := row.Schema(); v != schema {
		t.Errorf("incorrect schema: %v, expected: %v", v, schema)
	}

	if v := row.Metadata().UserID(); v != userID {
		t.Errorf("incorrect user id: %v, expected: %v", v, userID)
	}

	if v := row.Metadata().CorrelationID(); v != correlationID {
		t.Errorf("incorrect correlation id: %v, expected: %v", v, correlationID)
	}

	if v := row.Metadata().ClientVersion(); v != clientVersion {
		t.Errorf("incorrect client version: %v, expected: %v", v, clientVersion)
	}

	if v := row.Metadata().Time(); v != tm {
		t.Errorf("incorrect time: %v, expected: %v", v, tm)
	}

	if v := row.Data().Value(); !slices.Equal(v, data.Value()) {
		t.Errorf("incorrect name: %v, expected: %v", v, data.Value())
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
	cver := value.NewClientVersion("clientversion")
	tm := value.NewTime(time.Now())

	pointer := &Row{}

	for message, c := range map[string]struct {
		skip  bool
		a     *Row
		b     *Row
		equal bool
	}{
		"should equal zero": {
			a:     &Row{},
			b:     &Row{},
			equal: true,
		},
		"should equal pointer": {
			a:     pointer,
			b:     pointer,
			equal: true,
		},
		"should equal same values": {
			a: NewRow(
				rid,
				v,
				pid,
				rn,
				sch,
				NewMetadata(
					uid,
					cid,
					cver,
					tm,
				),
				value.NewData([]byte{1, 2, 3}),
			),
			b: NewRow(
				rid,
				v,
				pid,
				rn,
				sch,
				NewMetadata(
					uid,
					cid,
					cver,
					tm,
				),
				value.NewData([]byte{1, 2, 3}),
			),
			equal: true,
		},
		"should not equal different values": {
			a: NewRow(
				rid,
				v,
				pid,
				rn,
				sch,
				NewMetadata(
					uid,
					cid,
					cver,
					tm,
				),
				value.NewData([]byte{1, 2, 3}),
			),
			b: NewRow(
				value.NewRowID(),
				v,
				pid,
				rn,
				sch,
				NewMetadata(
					uid,
					cid,
					cver,
					tm,
				),
				value.NewData([]byte{1, 2, 3}),
			),
			equal: false,
		},
		"should not equal different data": {
			a: NewRow(
				rid,
				v,
				pid,
				rn,
				sch,
				NewMetadata(
					uid,
					cid,
					cver,
					tm,
				),
				value.NewData([]byte{1, 2, 3}),
			),
			b: NewRow(
				rid,
				v,
				pid,
				rn,
				sch,
				NewMetadata(
					uid,
					cid,
					cver,
					tm,
				),
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
