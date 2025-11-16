package wrap

import "encoding/binary"

type Int32 struct {
	data  []byte
	value int32
}

func NewInt32(data []byte) (Int32, []byte) {
	return Int32{data: data[:4]}, data[4:]
}

func (i *Int32) Get() int32 {
	return int32(binary.BigEndian.Uint32(i.data))
}

func (i *Int32) Set(v int32) {
	binary.BigEndian.PutUint32(i.data, uint32(v))
}
