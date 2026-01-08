package bytelist

type ByteList struct {
	data []byte
}

func New(data []byte) ByteList {
	return ByteList{
		data: data,
	}
}

func (l ByteList) Slice(off int64, n int64) ([]byte, bool) {
	if n == 0 {
		return l.data[off:off], true
	}

	if off >= int64(len(l.data)) {
		return nil, false
	}

	end := off + n
	if end > int64(len(l.data)) {
		return nil, false
	}

	return l.data[off:end], true
}
