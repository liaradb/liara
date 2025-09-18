package storage

import (
	"context"
	"io"

	"github.com/liaradb/liaradb/file"
)

type BufferManager struct {
	bufferSize int64
	fs         file.FS
	requests   chan *Buffer
}

// TODO: Should this be public?
func NewBufferManager(fs file.FS, bs int64) *BufferManager {
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

	_, err = f.ReadAt(b.data, bm.offset(b.blockID))
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (bm *BufferManager) Flush(b *Buffer) error {
	f, err := bm.openFile(b)
	if err != nil {
		return err
	}

	_, err = f.WriteAt(b.data, bm.offset(b.blockID))
	return err
}

func (bm *BufferManager) openFile(b *Buffer) (file.File, error) {
	return bm.fs.Open(b.blockID.FileName)
}

func (bm *BufferManager) offset(bid BlockID) int64 {
	return int64(bid.Offset(bm.bufferSize))
}
