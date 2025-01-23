package disk

import (
	"bytes"
	"fmt"
	"testing"
)

func newTestData() ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	rd := NewRecordEncoder(buffer)

	for _, value := range []string{
		"test0",
		"test1",
		"test2",
		"test3",
		"test4",
		"test5",
	} {
		_, _, err := rd.WriteRecord(value)
		if err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func TestRecordDecoder_Replay(t *testing.T) {
	data, err := newTestData()
	if err != nil {
		t.Fatal(err)
	}

	buffer := bytes.NewBuffer(data)

	rd := NewRecordDecoder(buffer)

	for record, err := range rd.Replay() {
		fmt.Print(record, err)
	}
}

func TestRecordDecoder_Get(t *testing.T) {
	data, err := newTestData()
	if err != nil {
		t.Fatal(err)
	}

	buffer := bytes.NewBuffer(data)

	rd := NewRecordDecoder(buffer)

	record, _ := rd.Get(2)
	if record.Value != "test2" {
		t.Error("wrote record")
	}
}

// func (d *Disk) replay() {
// 	d.file.ReadFrom()

// 	w := bufio.NewWriter(d.file)
// 	w.Write()
// 	w.Flush()

// 	r := bufio.NewReader(d.file)
// 	// r.ReadLine()

// 	decoder := json.NewDecoder(d.file)
// 	record := Record{}
// 	err := decoder.Decode(&record)
// 	// decoder.
// }
