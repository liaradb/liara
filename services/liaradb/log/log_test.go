package log

import "testing"

func TestLog_EmptyReader(t *testing.T) {
	t.Parallel()

	fsys := createFiles(0, 0)

	l := NewLog(256, fsys, ".")
	_, err := l.Reader(0)
	if err != ErrNoSegmentFile {
		t.Error("should have no files")
	}
}
