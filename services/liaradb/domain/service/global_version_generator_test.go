package service

import (
	"testing"

	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/domain/value"
)

func TestGlobalVersionGenerator(t *testing.T) {
	t.Parallel()

	g := GlobalVersionGenerator{}
	tn := tablename.NewFromString("tn")

	g.Init(tn, value.NewGlobalVersion(1))

	want := value.NewGlobalVersion(2)
	if v := g.Next(tn); v != want {
		t.Errorf("incorrect version: %v, expected: %v", v, want)
	}
}
