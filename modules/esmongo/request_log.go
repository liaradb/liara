package esmongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RequestLog struct {
	collection *Collection[request]
}

type RequestID string

func (i RequestID) String() string { return string(i) }

type request struct {
	ID   RequestID `bson:"_id"`
	Time time.Time `bson:"time"`
}

func NewRequestLog(
	collection *mongo.Collection,
) *RequestLog {
	return &RequestLog{
		collection: NewCollection[request](collection),
	}
}

func (r *RequestLog) Insert(ctx context.Context, requestID string, t time.Time) error {
	return r.collection.Insert(ctx, requestID, &request{
		ID:   RequestID(requestID),
		Time: t,
	})
}

func (r *RequestLog) Purge(ctx context.Context, t time.Time) error {
	f := bson.M{
		"time": bson.M{
			"lt": t}}
	_, err := r.collection.collection.DeleteMany(ctx, f)
	return err
}

func (r *RequestLog) Test(ctx context.Context, requestID string) (bool, error) {
	_, err := r.collection.Get(ctx, requestID)
	if errors.Is(err, ErrNotFound) {
		return true, nil
	}

	return false, err
}

func (r *RequestLog) Transaction(ctx context.Context, requestID string, t time.Time, tx func() error) (err error) {
	ok, err := r.Test(ctx, requestID)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	session, err := r.collection.collection.Database().Client().StartSession()
	if err != nil {
		return err
	}

	if err := session.StartTransaction(); err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, session.AbortTransaction(ctx))
		} else {
			err = session.CommitTransaction(ctx)
		}
	}()

	err = tx()
	if err != nil {
		return err
	}

	err = r.Insert(ctx, requestID, t)

	return
}
