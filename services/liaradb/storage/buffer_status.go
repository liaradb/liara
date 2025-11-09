package storage

type BufferStatus int

const (
	BufferStatusUninitialized BufferStatus = iota
	BufferStatusLoading
	BufferStatusLoaded
	BufferStatusDirty
	BufferStatusCorrupt
)
