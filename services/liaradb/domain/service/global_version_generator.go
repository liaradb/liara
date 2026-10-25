package service

import (
	"sync"

	"github.com/liaradb/liaradb/domain/value"
)

type GlobalVersionGenerator struct {
	mux      sync.Mutex
	versions map[value.TenantID]*tableVersion
}

func (g *GlobalVersionGenerator) Init(tid value.TenantID, v value.GlobalVersion) {
	tv := g.getTableVersion(tid)
	tv.setValue(v)
}

func (g *GlobalVersionGenerator) getTableVersion(tid value.TenantID) *tableVersion {
	g.mux.Lock()
	defer g.mux.Unlock()

	if g.versions == nil {
		g.versions = make(map[value.TenantID]*tableVersion)
	}

	tv, ok := g.versions[tid]
	if !ok {
		tv = &tableVersion{}
		g.versions[tid] = tv
	}
	return tv
}

func (g *GlobalVersionGenerator) Next(tid value.TenantID) value.GlobalVersion {
	tv := g.getTableVersion(tid)
	return tv.increment()
}
