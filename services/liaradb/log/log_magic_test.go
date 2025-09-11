package log

import (
	"io"
	"testing"
)

func TestLogMagic(t *testing.T) {
	r, w := createReaderWriter()

	var l LogMagic = LogMagicPage
	if err := l.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var l2 LogMagic
	if err := l2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if l != l2 {
		t.Errorf("incorrect value: %v, expected: %v", l2, l)
	}
}

func TestLogMagicPage(t *testing.T) {
	if s := LogMagicPage.String(); s != "PAGE" {
		t.Error("value is incorrect")
	}
}
