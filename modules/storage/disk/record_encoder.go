package disk

import (
	"bytes"
	"encoding/json"
	"io"
)

type RecordEncoder struct {
	index   uint64
	buffer  *bytes.Buffer
	encoder *json.Encoder
	writer  io.Writer
}

func NewRecordEncoder(
	writer io.Writer,
) *RecordEncoder {
	buffer := bytes.NewBuffer(nil)

	return &RecordEncoder{
		index:   0,
		buffer:  buffer,
		encoder: json.NewEncoder(buffer),
		writer:  writer,
	}
}

func (re *RecordEncoder) WriteRecord(value string) (uint64, int64, error) {
	re.buffer.Truncate(0)
	defer func() { re.index++ }()

	err := re.encoder.Encode(Record{Index: re.index, Value: value})
	if err != nil {
		return 0, 0, err
	}

	written, err := io.Copy(re.writer, re.buffer)
	return re.index, written, err
}
