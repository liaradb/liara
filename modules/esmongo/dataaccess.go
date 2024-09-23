package esmongo

import (
	"context"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DataAccess struct {
	client   *mongo.Client
	database *mongo.Database
}

func DbConnect(
	ctx context.Context,
	uri string,
	databaseName string,
	m *event.CommandMonitor,
) (*DataAccess, error) {
	dataAccess := new(DataAccess)
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

	dataAccess.client = client
	dataAccess.database = client.Database(databaseName)
	return dataAccess, nil
}

func (da *DataAccess) Database(databaseName string) DataAccess {
	return DataAccess{
		client:   da.client,
		database: da.client.Database(databaseName),
	}
}

func (da *DataAccess) Collection(collectionName string) *mongo.Collection {
	return da.database.Collection(collectionName)
}

func (da *DataAccess) Close(ctx context.Context) error {
	err := da.client.Disconnect(ctx)
	return err
}
