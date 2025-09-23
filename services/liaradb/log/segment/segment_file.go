package segment

import (
	"github.com/liaradb/liaradb/file"
)

type segmentFile struct {
	file file.File
	fsys file.FileSystem
	sn   SegmentName
}

func newSegmentFile(fsys file.FileSystem) *segmentFile {
	return &segmentFile{
		fsys: fsys,
	}
}

func (sf *segmentFile) SegmentName() SegmentName { return sf.sn }
func (sf *segmentFile) File() file.File          { return sf.file }

func (sf *segmentFile) Close() error {
	if !sf.isOpen() {
		return nil
	}

	if err := sf.closeFile(); err != nil {
		return err
	}

	return nil
}

func (sf *segmentFile) Open(sn SegmentName) (file.File, error) {
	// TODO: Test this
	if sf.isCurrentAndOpen(sn) {
		return sf.file, nil
	}

	if err := sf.Close(); err != nil {
		return nil, err
	}

	if err := sf.openFile(sn); err != nil {
		return nil, err
	}

	return sf.file, nil
}

func (sf *segmentFile) Remove(sn SegmentName) error {
	// TODO: Test this
	if sf.isCurrent(sn) {
		sf.Close()
	}

	if err := sf.fsys.Remove(sn.String()); err != nil {
		return err
	}

	return nil
}

func (sf *segmentFile) closeFile() error {
	if err := sf.file.Close(); err != nil {
		return err
	}

	sf.file = nil
	return nil
}

func (sf *segmentFile) isCurrent(sn SegmentName) bool {
	return sf.sn == sn
}

func (sf *segmentFile) isCurrentAndOpen(sn SegmentName) bool {
	return sf.isCurrent(sn) && sf.isOpen()
}

func (sf *segmentFile) isOpen() bool {
	return sf.file != nil
}

func (sf *segmentFile) openFile(sn SegmentName) error {
	f, err := sf.fsys.OpenFile(sn.String())
	if err != nil {
		return err
	}

	sf.sn = sn
	sf.file = f
	return nil
}
