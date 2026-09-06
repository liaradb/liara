package service

import (
	"sync"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/value"
)

type GlobalVersionGenerator struct {
	mux      sync.Mutex
	versions map[tablename.TableName]*tableVersion
}

func (g *GlobalVersionGenerator) Init(tn tablename.TableName, v value.GlobalVersion) {
	tv := g.getTableVersion(tn)
	tv.setValue(v)
}

func (g *GlobalVersionGenerator) getTableVersion(tn tablename.TableName) *tableVersion {
	g.mux.Lock()
	defer g.mux.Unlock()

	if g.versions == nil {
		g.versions = make(map[tablename.TableName]*tableVersion)
	}

	tv, ok := g.versions[tn]
	if !ok {
		tv = &tableVersion{}
		g.versions[tn] = tv
	}
	return tv
}

func (g *GlobalVersionGenerator) Next(tn tablename.TableName) value.GlobalVersion {
	tv := g.getTableVersion(tn)
	return tv.increment()
}
