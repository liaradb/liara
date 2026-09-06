package service

import (
	"sync"

	"github.com/liaradb/liaradb/domain/value"
)

type tableVersion struct {
	mux     sync.Mutex
	version value.GlobalVersion
}

func (tv *tableVersion) setValue(v value.GlobalVersion) {
	tv.mux.Lock()
	defer tv.mux.Unlock()

	tv.version = v
}

func (tv *tableVersion) increment() value.GlobalVersion {
	tv.mux.Lock()
	defer tv.mux.Unlock()

	tv.version = tv.version.Increment()
	return tv.version
}
