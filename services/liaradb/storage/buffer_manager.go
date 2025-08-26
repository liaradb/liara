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
	return newBuffer(bid, bm)
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

func (bm *BufferManager) openFile(b *Buffer) (file, error) {
	return bm.fs.Open(b.blockID.FileName)
}

func (bm *BufferManager) offset(bid BlockID) int64 {
	return int64(bid.Offset(bm.bufferSize))
}
