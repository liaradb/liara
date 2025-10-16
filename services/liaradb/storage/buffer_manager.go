package storage

import (
	"context"
	"io"

	"github.com/liaradb/liaradb/file"
)

type BufferManager struct {
	bufferSize int64
	fs         file.FileSystem
	requests   chan *Buffer
}

// TODO: Should this be public?
func NewBufferManager(fs file.FileSystem, bs int64) *BufferManager {
	return &BufferManager{
		bufferSize: bs,
		fs:         fs,
		requests:   make(chan *Buffer),
	}
}

func (bm *BufferManager) Request(ctx context.Context, b *Buffer) {
	select {
	case bm.requests <- b:
	case <-ctx.Done():
	}
}

func (bm *BufferManager) Load(b *Buffer) error {
	f, err := bm.openFile(b)
	if err != nil {
		return err
	}

	// TODO: Do we need to check io.EOF?
	if err := b.read(f); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (bm *BufferManager) Flush(b *Buffer) error {
	f, err := bm.openFile(b)
	if err != nil {
		return err
	}

	return b.write(f)
}

func (bm *BufferManager) openFile(b *Buffer) (file.File, error) {
	return bm.fs.OpenFile(b.blockID.FileName)
}
