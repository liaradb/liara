package raw

import "strconv"

type Offset int64

const (
	Uint32Length       Offset = 4
	Uint64Length       Offset = 8
	IntLength          Offset = strconv.IntSize / 8
	Int32Length        Offset = 4
	Int64Length        Offset = 8
	stringHeaderOffset Offset = Uint32Length
)

func (o Offset) Value() int64 { return int64(o) }

func StrSizeFromLength(length int) Offset {
	return Uint32Length + Offset(length)
}

func StrSize(value string) Offset {
	return Uint32Length + Offset(len(value))
}

func BufferSize(bytes []byte) Offset {
	return Uint32Length + Offset(len(bytes))
}
