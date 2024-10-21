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
var ErrNotFoundAtVersion = errors.New("not found at version")

type Page struct {
	Offset int
	Limit  int
}

type Mapper[T any, M any] interface {
	FromModel(*M) *T
	ToModel(*T) *M
}

type MongoRepository[T any, I ~string, M any] struct {
	collection *mongo.Collection
	mapper     Mapper[T, M]
}

func NewMongoRepository[T any, I ~string, M any](
	collection *mongo.Collection,
	mapper Mapper[T, M],
) *MongoRepository[T, I, M] {
	return &MongoRepository[T, I, M]{
		collection: collection,
		mapper:     mapper,
	}
}

func (mr *MongoRepository[T, I, M]) Insert(
	ctx context.Context,
	id I,
	t *T,
) error {
	_, err := mr.collection.ReplaceOne(ctx,
		bson.M{"_id": id},
		mr.mapper.ToModel(t),
		options.Replace().SetUpsert(true))
	return err
}

func (mr *MongoRepository[T, I, M]) Replace(
	ctx context.Context,
	id I,
	v int,
	t *T,
) error {
	result, err := mr.collection.ReplaceOne(ctx,
		bson.M{
			"_id":     id,
			"version": v},
		mr.mapper.ToModel(t),
		options.Replace())
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrNotFoundAtVersion
	}

	return nil
}

func (mr *MongoRepository[T, I, M]) Get(
	ctx context.Context,
	id I,
) (*T, error) {
	m, err := decode[M](mr.collection.FindOne(ctx, bson.M{
		"_id": id}))
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = ErrNotFound
	}

	e := mr.mapper.FromModel(&m)
	return e, err
}

func (mr *MongoRepository[T, I, M]) Delete(
	ctx context.Context,
	id I,
) error {
	_, err := mr.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (mr *MongoRepository[T, I, M]) GetList(
	ctx context.Context,
	filter FilterBuilder,
	sort *SortBuilder,
) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		f := filter.Build()
		o := sort.Build()
		result, err := mr.collection.Find(ctx, f, o)
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

			e := mr.mapper.FromModel(&m)
			if yield(e, nil) {
				return
			}
		}
	}
}

func (mr *MongoRepository[T, I, M]) Watch(
	ctx context.Context,
	pipeline any,
	token string,
) (iter.Seq2[bson.Raw, string], error) {
	o := options.ChangeStream()
	if token != "" {
		o.SetResumeAfter(bson.Raw(token))
	}

	cs, err := mr.collection.Watch(ctx, pipeline, o)
	if err != nil {
		return nil, err
	}

	return func(yield func(bson.Raw, string) bool) {
		defer cs.Close(ctx)

		for cs.Next(ctx) {
			if !yield(cs.Current, cs.ResumeToken().String()) {
				return
			}
		}
	}, nil
}

func RunTransaction[T any](
	ctx context.Context,
	c *mongo.Client,
	p func(ctx context.Context) (T, error),
) (T, error) {
	s, err := c.StartSession()
	if err != nil {
		var t T
		return t, err
	}
	defer s.EndSession(ctx)

	value, err := s.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		return p(ctx)
	})
	t, _ := value.(T)
	return t, err
}
