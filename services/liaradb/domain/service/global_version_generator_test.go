package service

import (
	"testing"

	"github.com/liaradb/liaradb/domain/value"
)

func TestGlobalVersionGenerator(t *testing.T) {
	t.Parallel()

	g := GlobalVersionGenerator{}
	tid := value.NewTenantID()

	g.Init(tid, value.NewGlobalVersion(1))

	want := value.NewGlobalVersion(2)
	if v := g.Next(tid); v != want {
		t.Errorf("incorrect version: %v, expected: %v", v, want)
	}
}
