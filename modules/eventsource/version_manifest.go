package eventsource

import "github.com/cardboardrobots/eventsource/value"

type VersionManifest struct {
	versions map[value.AggregateID]value.Version
}

func (vm *VersionManifest) AddVersion(id value.AggregateID, version value.Version) bool {
	if version <= vm.versions[id] {
		return true
	}

	if vm.versions == nil {
		vm.versions = make(map[value.AggregateID]value.Version)
	}

	vm.versions[id] = version

	return false
}

func (vm *VersionManifest) GetVersion(id value.AggregateID) value.Version {
	return vm.versions[id]
}
