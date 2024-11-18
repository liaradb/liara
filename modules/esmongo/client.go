package esmongo

import (
	"context"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect(
	ctx context.Context,
	uri string,
	m *event.CommandMonitor,
) (*mongo.Client, error) {
	o := options.Client().ApplyURI(uri)
	if m != nil {
		o = o.SetMonitor(m)
	}
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
