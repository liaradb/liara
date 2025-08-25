package storage

import "io"

type BufferManager struct {
	bufferSize int64
	fs         FS
}

func newBufferManager(fs FS) *BufferManager {
	return &BufferManager{
		bufferSize: 1024,
		fs:         fs,
	}
}

func (bm *BufferManager) Buffer(bid BlockID) *Buffer {
	return &Buffer{
		blockID: bid,
		data:    make([]byte, bm.bufferSize),
		bm:      bm,
	}
}

func (bm *BufferManager) Load(b *Buffer) error {
	f, err := bm.fs.Open(b.blockID.FileName)
	if err != nil {
		return err
	}

	_, err = f.ReadAt(b.data, int64(b.blockID.Position)*bm.bufferSize)
	if err == io.EOF {
		return nil
	}
	return err
}

func (bm *BufferManager) Flush(b *Buffer) error {
	f, err := bm.fs.Open(b.blockID.FileName)
	if err != nil {
		return err
	}

	_, err = f.WriteAt(b.data, int64(b.blockID.Position)*bm.bufferSize)
	return err
}
