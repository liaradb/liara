package eventlog

import (
	"context"

	"github.com/liaradb/liaradb/storage"
)

type Tail struct {
	storage *storage.Storage
}

func NewTail(
	storage *storage.Storage,
) *Tail {
	return &Tail{
		storage: storage,
	}
}

func (t *Tail) Append(ctx context.Context) error {
	return nil
}
