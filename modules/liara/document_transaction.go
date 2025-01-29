package liara

import (
	"context"
	"errors"
	"time"
)

type DocumentTransaction struct {
	requestLog RequestLog
	handler    EventHandler
}

type RequestLog interface {
	Test(context.Context, string) (bool, error)
	Insert(context.Context, string, time.Time) error
	Purge(context.Context, time.Time) error
	Transaction(context.Context, string, time.Time, func() error) error
}

func NewDocumentTransaction(
	requestLog RequestLog,
	handler EventHandler,
) *DocumentTransaction {
	return &DocumentTransaction{
		requestLog: requestLog,
		handler:    handler,
	}
}

func (l *DocumentTransaction) Handle(ctx context.Context, em Event) error {
	return l.requestLog.Transaction(ctx, em.ID.String(), time.Now(), func() error {
		err := l.handler.Handle(ctx, em)

		// Skip missing streams
		if errors.Is(err, ErrNotFound) {
			return nil
		}

		return err
	})
}
