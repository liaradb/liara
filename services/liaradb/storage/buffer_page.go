package storage

type BufferPage struct {
	*Buffer
}

func NewBufferPage(b *Buffer) *BufferPage {
	return &BufferPage{
		Buffer: b,
	}
}
