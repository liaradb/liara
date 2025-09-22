package segment

import (
	"github.com/liaradb/liaradb/file"
	"github.com/liaradb/liaradb/log/record"
)

type segmentFile struct {
	sn   SegmentName
	fsys file.FileSystem
	file file.File
}

func newSegmentFile(sn SegmentName, fsys file.FileSystem) *segmentFile {
	return &segmentFile{
		sn:   sn,
		fsys: fsys,
	}
}

func (sf *segmentFile) SegmentName() SegmentName { return sf.sn }
func (sf *segmentFile) File() file.File          { return sf.file }

func (sf *segmentFile) Close() error {
	if sf.file == nil {
		return nil
	}

	if err := sf.file.Close(); err != nil {
		return err
	}

	sf.file = nil
	return nil
}

func (sf *segmentFile) Next(lsn record.LogSequenceNumber) *segmentFile {
	return newSegmentFile(sf.sn.Next(lsn), sf.fsys)
}

func (sf *segmentFile) Open() error {
	if sf.file != nil {
		return nil
	}

	f, err := sf.fsys.OpenFile(sf.sn.String())
	if err != nil {
		return err
	}

	sf.file = f
	return nil
}
