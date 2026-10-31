package lrupool

import (
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
	"github.com/liaradb/liaradb/util/queue"
)

func New() *queue.MapQueue[link.FileBlockID, *storage.Buffer] {
	return &queue.MapQueue[link.FileBlockID, *storage.Buffer]{}
}
