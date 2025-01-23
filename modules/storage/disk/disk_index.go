package disk

import (
	"encoding/json"
	"io"
)

type DiskIndex struct {
	encoder *json.Encoder
	decoder *json.Decoder
	data    map[uint64]uint64
}

type recordIndex struct {
	Index    uint64 `json:"index"`
	Position uint64 `json:"position"`
}

func NewDiskIndex(writer io.Writer, reader io.Reader) *DiskIndex {
	return &DiskIndex{
		encoder: json.NewEncoder(writer),
		decoder: json.NewDecoder(reader),
	}
}

func (di *DiskIndex) WriteIndex(
	index uint64,
	position uint64,
) error {
	err := di.encoder.Encode(recordIndex{
		Index:    index,
		Position: position,
	})
	if err != nil {
		return err
	}

	if di.data == nil {
		di.data = make(map[uint64]uint64)
	}

	di.data[index] = position

	return nil
}

func (di *DiskIndex) Restore() error {
	if di.data == nil {
		di.data = make(map[uint64]uint64)
	}

	ri := recordIndex{}
	for {
		err := di.decoder.Decode(&ri)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		di.data[ri.Index] = ri.Position
	}

	return nil
}

func (di *DiskIndex) Get(index uint64) (uint64, bool) {
	position, ok := di.data[index]
	return position, ok
}
