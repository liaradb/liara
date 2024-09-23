package esmongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Cast[T ~string, U ~string](a T) U {
	return U(a)
}

func PageToOptions(limit int64, offset int64) *options.FindOptions {
	o := options.Find()
	if limit > 0 {
		o.SetLimit(limit)
	}
	if offset > 0 {
		o.SetSkip(offset)
	}
	return o
}

func CursorToSlice[T any](
	ctx context.Context,
	c *mongo.Cursor,
) ([]T, error) {
	var results []T
	if err := c.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func CursorToEntitySlice[T any, U any](
	ctx context.Context,
	c *mongo.Cursor,
	m func(a T) U,
) ([]U, error) {
	var results []T
	if err := c.All(ctx, &results); err != nil {
		return nil, err
	}
	u := make([]U, 0, len(results))
	for _, r := range results {
		u = append(u, m(r))
	}
	return u, nil
}

func FindToEntitySlice[T any, U any](
	ctx context.Context,
	m func(a T) U,
	c *mongo.Collection,
	f any,
	o ...*options.FindOptions,
) ([]U, error) {
	cursor, err := c.Find(ctx, f, o...)
	if err != nil {
		return nil, err
	}
	return CursorToEntitySlice[T, U](ctx, cursor, m)
}

func FindEntityList[T any, U any](
	ctx context.Context,
	m func(a T) U,
	c *mongo.Collection,
	f any,
	o ...*options.FindOptions,
) ([]U, int64, error) {
	cursor, err := c.Find(ctx, f, o...)
	if err != nil {
		return nil, 0, err
	}

	data, err := CursorToEntitySlice[T, U](ctx, cursor, m)
	if err != nil {
		return nil, 0, err
	}

	count, err := c.CountDocuments(ctx, f)
	if err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func Count(
	ctx context.Context,
	c *mongo.Collection,
	f any,
	o ...*options.FindOptions,
) (int64, error) {
	count, err := c.CountDocuments(ctx, f)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func FindOne[T any, U any](
	ctx context.Context,
	c *mongo.Collection,
	m func(a T) U,
	filter any,
) (*U, error) {
	var result T
	if err := c.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	u := m(result)
	return &u, nil
}

func FindByID[T any, U any](
	ctx context.Context,
	c *mongo.Collection,
	m func(a T) U,
	id string,
) (*U, error) {
	var result T
	if err := c.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&result); err != nil {
		return nil, err
	}

	u := m(result)
	return &u, nil
}

func Create[T any, U any](
	ctx context.Context,
	c *mongo.Collection,
	m func(a T) U,
	a T,
) (string, error) {
	result, error := c.InsertOne(ctx, m(a))
	id := result.InsertedID.(string)
	return id, error
}

func ReplaceByID[T any, U any, V ~string](
	ctx context.Context,
	c *mongo.Collection,
	m func(a T) U,
	id V,
	a T,
) error {
	_, err := c.ReplaceOne(ctx, bson.M{
		"_id": id,
	}, m(a))
	return err
}

func DeleteByID[V ~string](
	ctx context.Context,
	c *mongo.Collection,
	id V,
) error {
	if _, err := c.DeleteOne(ctx, bson.M{
		"_id": id,
	}); err != nil {
		return err
	}
	return nil
}
