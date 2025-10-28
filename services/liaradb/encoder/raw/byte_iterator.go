package raw

type ByteIterator struct {
	bytes    []byte
	position Offset
}

func NewByteIterator(bytes []byte) ByteIterator {
	return ByteIterator{
		bytes: bytes,
	}
}

func (bi *ByteIterator) Bytes() []byte {
	return bi.bytes
}

func (bi *ByteIterator) Position() Offset {
	return bi.position
}

func (bi *ByteIterator) Reset() {
	bi.position = 0
}

func (bi *ByteIterator) GetUint32() (uint32, error) {
	i, err := GetUint32(bi.bytes, bi.position)
	if err != nil {
		return 0, err
	}

	bi.position += Uint32Length

	return i, nil
}

func (bi *ByteIterator) SetUint32(value uint32) error {
	if err := CopyUint32(bi.bytes, value, bi.position); err != nil {
		return err
	}

	bi.position += Uint32Length

	return nil
}

func (bi *ByteIterator) ForwardUint32() error {
	p := bi.position + Uint32Length
	if p >= Offset(len(bi.bytes)) {
		return ErrOverflow
	}

	bi.position = p
	return nil
}

func (bi *ByteIterator) BackUint32() {
	bi.position -= Uint32Length
	if bi.position < 0 {
		bi.position = 0
	}
}

func (bi *ByteIterator) GetString() (string, error) {
	s, err := GetString(bi.bytes, bi.position)
	if err != nil {
		return "", err
	}

	bi.position += Uint32Length + Offset(len(s))

	return s, nil
}

func (bi *ByteIterator) SetString(value string) error {
	if err := CopyString(bi.bytes, value, bi.position); err != nil {
		return err
	}

	bi.position += Uint32Length + Offset(len(value))

	return nil
}

func (bi *ByteIterator) ForwardString(value string) error {
	p := bi.position + Uint32Length + Offset(len(value))
	if p >= Offset(len(bi.bytes)) {
		return ErrOverflow
	}

	bi.position = p
	return nil
}

func (bi *ByteIterator) InitBytes() {
	bi.bytes = make([]byte, bi.position+1)
}
