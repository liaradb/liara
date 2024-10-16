package eventsource

import "github.com/cardboardrobots/liara"

type VersionManifest struct {
	versions map[liara.AggregateID]liara.Version
}

func (vm *VersionManifest) AddVersion(id liara.AggregateID, version liara.Version) bool {
	if version <= vm.versions[id] {
		return true
	}

	if vm.versions == nil {
		vm.versions = make(map[liara.AggregateID]liara.Version)
	}

	vm.versions[id] = version

	return false
}

func (vm *VersionManifest) GetVersion(id liara.AggregateID) liara.Version {
	return vm.versions[id]
}
