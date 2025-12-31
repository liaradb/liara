package link

type FileName string

func NewFileName(value string) FileName {
	return FileName(value)
}

func (fn FileName) String() string { return string(fn) }

func (fn FileName) BlockID(position FilePosition) BlockID {
	return NewBlockID(fn, position)
}
