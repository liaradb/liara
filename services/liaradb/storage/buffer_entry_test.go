package storage

import "testing"

func TestBufferEntry(t *testing.T) {
	b, close := testCreateBuffer(t)
	defer close()

	if err := b.Load(); err != nil {
		t.Fatal(err)
	}

	number := newUInt64Entry(0)

	var want uint64 = 12345

	if err := number.Set(b, want); err != nil {
		t.Fatal(err)
	}

	if v, err := number.Get(b); err != nil {
		t.Error(err)
	} else if v != want {
		t.Errorf("value does not match: expected: %v, recieved: %v", want, v)
	}
}
