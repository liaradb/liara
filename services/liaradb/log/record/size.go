package record

type sizer interface {
	Size() int
}

func size(sizers ...sizer) int {
	size := 0
	for _, s := range sizers {
		size += s.Size()
	}
	return size
}
