package slice

func Slice(data []byte, off int64, n int64) ([]byte, bool) {
	if n == 0 {
		return data[off:off], true
	}

	if off >= int64(len(data)) {
		return nil, false
	}

	end := off + n
	if end > int64(len(data)) {
		return nil, false
	}

	return data[off:end], true
}
