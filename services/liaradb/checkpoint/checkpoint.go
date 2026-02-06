package checkpoint

import (
	"time"

	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/storage"
)

type Checkpoint struct {
	l *recovery.Log
	s *storage.Storage
}

func New(
	l *recovery.Log,
	s *storage.Storage,
) *Checkpoint {
	return &Checkpoint{
		l: l,
		s: s,
	}
}

func (c *Checkpoint) Flush(now time.Time) error {
	if err := c.s.FlushAll(); err != nil {
		return err
	}

	_, err := c.l.FlushCheckpoint(now)
	return err
}
