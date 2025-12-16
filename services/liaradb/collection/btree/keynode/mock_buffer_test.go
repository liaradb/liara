package keynode

type mockBuffer struct {
	data []byte
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func newMockBuffer[T Number](size T) *mockBuffer {
	return &mockBuffer{
		data: make([]byte, size),
	}
}

func (b *mockBuffer) Raw() []byte { return b.data }
func (b *mockBuffer) Clear()      { clear(b.data) }
func (b *mockBuffer) Release()    {}
func (b *mockBuffer) SetDirty()   {}
