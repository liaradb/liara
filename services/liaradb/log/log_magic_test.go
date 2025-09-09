package log

import "testing"

func TestLogMagicPage(t *testing.T) {
	if s := LogMagicPage.String(); s != "PAGE" {
		t.Error("value is incorrect")
	}
}
