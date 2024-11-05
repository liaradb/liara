package esmongo

import (
	"context"
	"errors"
	"iter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("not found")

type Collection[M any] struct {
	collection *mongo.Collection
}

func NewCollection[M any](collection *mongo.Collection) *Collection[M] {
	return &Collection[M]{collection}
}

func (c *Collection[M]) Insert(ctx context.Context, id string, m *M) error {
	_, err := c.collection.ReplaceOne(ctx,
		bson.M{"_id": id},
		m)
	return err
}

func (c *Collection[M]) Upsert(ctx context.Context, id string, m *M) error {
	_, err := c.collection.ReplaceOne(ctx,
		bson.M{"_id": id},
		m,
		options.Replace().SetUpsert(true))
	return err
}

func (c *Collection[M]) Replace(ctx context.Context, filter FilterBuilder, m *M) error {
	result, err := c.collection.ReplaceOne(ctx,
		filter.Build(),
		m)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (c *Collection[M]) Get(ctx context.Context, id string) (*M, error) {
	m, err := decode[M](c.collection.FindOne(ctx,
		bson.M{"_id": id}))
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = ErrNotFound
	}

	return &m, err
}

func (c *Collection[M]) Delete(ctx context.Context, id string) error {
	_, err := c.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (c *Collection[M]) GetList(ctx context.Context, filter FilterBuilder, sort *SortBuilder) iter.Seq2[*M, error] {
	return func(yield func(*M, error) bool) {
		f := filter.Build()
		o := sort.Build()
		result, err := c.collection.Find(ctx, f, o)
		if err != nil {
			yield(nil, err)
			return
		}

		defer result.Close(ctx)

		for result.Next(ctx) {
			m, err := decode[M](result)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					err = ErrNotFound
				}
				yield(nil, err)
				return
			}

			if yield(&m, nil) {
				return
			}
		}
	}
}

type Change[M any] struct {
	Value M
	Token string
}

func (c *Collection[M]) Watch(ctx context.Context, pipeline any, token string) iter.Seq2[Change[M], error] {
	return func(yield func(Change[M], error) bool) {
		o := options.ChangeStream()
		if token != "" {
			o.SetResumeAfter(bson.Raw(token))
		}

		cs, err := c.collection.Watch(ctx, pipeline, o)
		if err != nil {
			yield(Change[M]{Token: token}, err)
			return
		}

		defer cs.Close(ctx)

		for cs.Next(ctx) {
			m, err := decode[M](cs)
			if err != nil {
				yield(Change[M]{Token: cs.ResumeToken().String()}, nil)
				return
			}

			if !yield(Change[M]{Value: m, Token: cs.ResumeToken().String()}, nil) {
				return
			}
		}
	}
}
