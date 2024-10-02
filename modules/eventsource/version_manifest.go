package eventsource

type VersionManifest struct {
	versions map[AggregateID]Version
}

func (vm *VersionManifest) AddVersion(id AggregateID, version Version) bool {
	if version <= vm.versions[id] {
		return true
	}

	if vm.versions == nil {
		vm.versions = make(map[AggregateID]Version)
	}

	vm.versions[id] = version

	return false
}

func (vm *VersionManifest) GetVersion(id AggregateID) Version {
	return vm.versions[id]
}
