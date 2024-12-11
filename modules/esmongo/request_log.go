package esmongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type RequestLogInterface interface {
	Test(context.Context, TenantID, RequestID) (bool, error)
	Insert(context.Context, TenantID, RequestID, time.Time) error
	Purge(context.Context, TenantID, time.Time) error
}

type RequestLog struct {
	collection Collection[request]
}

var _ RequestLogInterface = (*RequestLog)(nil)

type (
	RequestID string
	TenantID  string
)

func (i RequestID) String() string { return string(i) }
func (i TenantID) String() string  { return string(i) }

type request struct {
	ID   RequestID `bson:"_id"`
	Time time.Time `bson:"time"`
}

func (r *RequestLog) Insert(ctx context.Context, tenantID TenantID, requestID RequestID, t time.Time) error {
	return r.collection.Insert(ctx, requestID.String(), &request{
		ID:   requestID,
		Time: t,
	})
}

func (r *RequestLog) Purge(ctx context.Context, tenantID TenantID, t time.Time) error {
	f := bson.M{
		"time": bson.M{
			"lt": t}}
	_, err := r.collection.collection.DeleteMany(ctx, f)
	return err
}

func (r *RequestLog) Test(ctx context.Context, tenantID TenantID, requestID RequestID) (bool, error) {
	_, err := r.collection.Get(ctx, requestID.String())
	if errors.Is(err, ErrNotFound) {
		return true, nil
	}

	return false, err
}
