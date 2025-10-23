package value

type Version struct {
	baseUint64
}

func NewVersion(value uint64) Version {
	return Version{baseUint64(value)}
}
