package wrap

import "encoding/binary"

type Int64 struct {
	data  []byte
	value int64
}

func NewInt64(data []byte) (Int64, []byte) {
	return Int64{data: data[:8]}, data[8:]
}

func (i *Int64) Get() int64 {
	return int64(binary.BigEndian.Uint64(i.data))
}

func (i *Int64) GetUnsigned() uint64 {
	return binary.BigEndian.Uint64(i.data)
}

func (i *Int64) Set(v int64) {
	binary.BigEndian.PutUint64(i.data, uint64(v))
}

func (i *Int64) SetUnsigned(v uint64) {
	binary.BigEndian.PutUint64(i.data, v)
}
