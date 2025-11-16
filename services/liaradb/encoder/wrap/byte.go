package wrap

type Byte struct {
	data  []byte
	value byte
}

func NewByte(data []byte) (Byte, []byte) {
	return Byte{data: data[:1]}, data[1:]
}

func (i *Byte) Get() int8 {
	return int8(i.data[0])
}

func (i *Byte) GetUnsigned() byte {
	return i.data[0]
}

func (i *Byte) Set(v int8) {
	i.data[0] = byte(v)
}

func (i *Byte) SetUnsigned(v byte) {
	i.data[0] = v
}
