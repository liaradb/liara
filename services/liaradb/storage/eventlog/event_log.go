package eventlog

import (
	"bufio"
	"bytes"
	"context"

	"github.com/liaradb/liaradb/storage"
)

type EventLog struct {
	storage *storage.Storage
	buffer  *bytes.Buffer
	reader  *bufio.Reader
}

func New(
	storage *storage.Storage,
) *EventLog {
	buffer := bytes.NewBuffer(nil)
	reader := bufio.NewReader(buffer)
	return &EventLog{
		storage: storage,
		buffer:  buffer,
		reader:  reader,
	}
}

func (l *EventLog) Append(ctx context.Context, fileName string, e *Event) error {
	var data []byte
	if _, err := l.buffer.Write(data); err != nil {
		return err
	}

	_, err := l.storage.Append(ctx, fileName, l.reader)
	if err != nil {
		return err
	}

	l.buffer.Reset()
	return nil
}
