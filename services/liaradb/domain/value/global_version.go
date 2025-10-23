package value

type GlobalVersion struct {
	baseUint64
}

func NewGlobalVersion(value uint64) GlobalVersion {
	return GlobalVersion{baseUint64(value)}
}
