package segment

import (
	"path"

	"github.com/liaradb/liaradb/file"
)

type segmentFile struct {
	file file.File
	fsys file.FileSystem
	dir  string
	sn   SegmentName
}

func newSegmentFile(fsys file.FileSystem, dir string) *segmentFile {
	return &segmentFile{
		fsys: fsys,
		dir:  dir,
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

func (sf *segmentFile) open(sn SegmentName) (file.File, error) {
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

func (sf *segmentFile) remove(sn SegmentName) error {
	if sf.isCurrent(sn) {
		if err := sf.Close(); err != nil {
			return err
		}
	}

	if err := sf.fsys.Remove(sf.path(sn)); err != nil {
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
	f, err := sf.fsys.OpenFile(sf.path(sn))
	if err != nil {
		return err
	}

	sf.sn = sn
	sf.file = f
	return nil
}

func (sf *segmentFile) path(sn SegmentName) string {
	return path.Join(sf.dir, sn.String())
}
