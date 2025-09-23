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
	if sf.file == nil {
		return nil
	}

	if err := sf.file.Close(); err != nil {
		return err
	}

	sf.file = nil
	return nil
}

func (sf *segmentFile) Open(sn SegmentName) (file.File, error) {
	// TODO: Test this
	if sf.sn == sn && sf.file != nil {
		return sf.file, nil
	}

	if err := sf.Close(); err != nil {
		return nil, err
	}

	sf.reset(sn)

	f, err := sf.fsys.OpenFile(sf.sn.String())
	if err != nil {
		return nil, err
	}

	sf.file = f
	return f, nil
}

func (sf *segmentFile) Remove(sn SegmentName) error {
	// TODO: Test this
	if sn == sf.sn {
		sf.Close()
	}

	if err := sf.fsys.Remove(sn.String()); err != nil {
		return err
	}

	return nil
}

func (sf *segmentFile) reset(sn SegmentName) {
	sf.sn = sn
	sf.file = nil

}
