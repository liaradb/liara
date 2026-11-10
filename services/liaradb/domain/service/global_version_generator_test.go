package service

import (
	"testing"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/transaction/record"
)

func TestGlobalVersionGenerator(t *testing.T) {
	t.Parallel()

	g := GlobalVersionGenerator{}
	tid := value.NewTenantID()
	txid := record.NewTransactionID(3)

	g.Init(tid, value.NewGlobalVersion(1))

	want := value.NewGlobalVersion(2)
	if v := g.Next(tid, txid); v != want {
		t.Errorf("incorrect version: %v, expected: %v", v, want)
	}
}
