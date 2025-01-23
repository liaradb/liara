package disk

import (
	"encoding/json"
	"errors"
	"io"
	"iter"
)

type RecordDecoder struct {
	decoder *json.Decoder
}

func NewRecordDecoder(
	reader io.Reader,
) *RecordDecoder {
	return &RecordDecoder{
		decoder: json.NewDecoder(reader),
	}
}

func (rd *RecordDecoder) Replay() iter.Seq2[Record, error] {
	return func(yield func(Record, error) bool) {
		record := Record{}
		for {
			err := rd.decoder.Decode(&record)
			if err == io.EOF {
				return
			}

			if !yield(record, err) {
				return
			}
		}
	}
}

// Gets the Record with the specified index
// Without indexing, this must replay from the beginning
func (rd *RecordDecoder) Get(index uint64) (Record, error) {
	for record, err := range rd.Replay() {
		if err != nil {
			return Record{}, err
		}

		if record.Index == index {
			return record, nil
		}
		if record.Index > index {
			break
		}
	}

	return Record{}, errors.New("not found")
}

func (rd *RecordDecoder) Decode() (Record, error) {
	record := Record{}
	err := rd.decoder.Decode(&record)
	return record, err
}
