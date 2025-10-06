package eventlog

import (
	"context"

	"github.com/liaradb/liaradb/storage"
)

type EventLog struct {
	storage *storage.Storage
}

func New(
	storage *storage.Storage,
) *EventLog {
	return &EventLog{
		storage: storage,
	}
}

func (l *EventLog) Append(ctx context.Context, fileName string, e *Event) error {
	_, err := l.storage.RequestLatest(ctx, fileName)
	if err != nil {
		return err
	}

	return nil
}
