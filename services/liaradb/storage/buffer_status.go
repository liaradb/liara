package storage

type BufferStatus int

const (
	BufferStatusUninitialized BufferStatus = iota
	BufferStatusLoading
	BufferStatusLoaded
	BufferStatusDirty
	BufferStatusCorrupt
)

func (bs BufferStatus) String() string {
	switch bs {
	case BufferStatusCorrupt:
		return "corrupt"
	case BufferStatusDirty:
		return "dirty"
	case BufferStatusLoaded:
		return "loaded"
	case BufferStatusLoading:
		return "loading"
	case BufferStatusUninitialized:
		return "uninitialized"
	default:
		return "unknown"
	}
}
