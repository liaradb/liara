package raw

type Sizer interface {
	Size() int
}

func Size(sizers ...Sizer) int {
	size := 0
	for _, s := range sizers {
		size += s.Size()
	}
	return size
}
