package eventsource

import "testing"

func TestVersionManifest_AddVersion(t *testing.T) {
	vm := VersionManifest{}
	vm.AddVersion("1", 1)
	version := vm.GetVersion("1")
	if version != 1 {
		t.Error("version mismatch")
	}
}
