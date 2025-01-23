package disk

import (
	"errors"
	"io"
	"iter"
	"os"
)

type Disk struct {
	pageSize int
	page     int
	offset   uint64
	file     *os.File
	encoder  *RecordEncoder
	decoder  *RecordDecoder
	index    *DiskIndex
}

func NewDisk() Disk {
	return Disk{
		pageSize: os.Getpagesize(),
		page:     0,
	}
}

func (d *Disk) Open(name string, indexName string) error {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	position, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	d.file = file
	d.offset = uint64(position)
	d.encoder = NewRecordEncoder(file)
	d.decoder = NewRecordDecoder(file)

	indexFile, err := os.OpenFile(indexName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	d.index = NewDiskIndex(indexFile, indexFile)
	return d.index.Restore()
}

func (d *Disk) Page() int { return d.page }

func (d *Disk) Seek(page int) {
	d.page = page
}

func (d *Disk) position() int64 {
	return int64(d.pageSize) * int64(d.page)
}

func (d *Disk) WriteRecord(value string) (uint64, error) {
	_, err := d.file.Seek(int64(d.offset), io.SeekStart)
	if err != nil {
		return 0, err
	}

	index, position, err := d.encoder.WriteRecord(value)
	if err != nil {
		return 0, err
	}

	err = d.index.WriteIndex(index, d.offset)
	// d.file.Sync()

	d.offset += uint64(position)

	return index, err
}

func (d *Disk) Replay() iter.Seq2[Record, error] {
	_, err := d.file.Seek(d.position(), 0)
	if err != nil {
		return func(yield func(Record, error) bool) {
			yield(Record{}, err)
		}
	}
	return d.decoder.Replay()
}

func (d *Disk) Get(index uint64) (Record, error) {
	position, ok := d.index.Get(index)
	if !ok {
		return Record{}, errors.New("not found")
	}

	_, err := d.file.Seek(int64(position), io.SeekStart)
	if err != nil {
		return Record{}, err
	}

	return d.decoder.Decode()
}
