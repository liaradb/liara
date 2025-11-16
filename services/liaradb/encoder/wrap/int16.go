package wrap

import "encoding/binary"

type Int16 struct {
	data  []byte
	value int16
}

func NewInt16(data []byte) (Int16, []byte) {
	return Int16{data: data[:2]}, data[2:]
}

func (i *Int16) Get() int16 {
	return int16(binary.BigEndian.Uint16(i.data))
}

func (i *Int16) GetUnsigned() uint16 {
	return binary.BigEndian.Uint16(i.data)
}

func (i *Int16) Set(v int16) {
	binary.BigEndian.PutUint16(i.data, uint16(v))
}

func (i *Int16) SetUnsigned(v uint16) {
	binary.BigEndian.PutUint16(i.data, v)
}
