package esmongo

import (
	"context"
	"errors"
	"iter"

	"github.com/cardboardrobots/liara"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Page struct {
	Offset int
	Limit  int
}

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI(uri))
}

func Database(c *mongo.Client, dbName string) *mongo.Database {
	return c.Database(dbName)
}

type decoder interface {
	Decode(any) error
}

func decode[M any](decoder decoder) (M, error) {
	var m M
	return m, decoder.Decode(&m)
}

type Mapper[T any, I ~string, M any] interface {
	FromM(*M) (liara.Version, *T)
	ToM(liara.Version, I, *T) *M
}

type MongoRepository[T any, I ~string, M any] struct {
	collection *mongo.Collection
	mapper     Mapper[T, I, M]
}

func NewMongoRepository[T any, I ~string, M any](
	collection *mongo.Collection,
	mapper Mapper[T, I, M],
) *MongoRepository[T, I, M] {
	return &MongoRepository[T, I, M]{
		collection: collection,
		mapper:     mapper,
	}
}

func (mr *MongoRepository[T, I, M]) Insert(
	ctx context.Context,
	v liara.Version,
	id I,
	t *T,
) error {
	_, err := mr.collection.ReplaceOne(ctx,
		bson.M{"_id": id},
		mr.mapper.ToM(v, id, t),
		options.Replace().SetUpsert(true))
	return err
}

func (mr *MongoRepository[T, I, M]) Get(
	ctx context.Context,
	id I,
) (liara.Version, *T, error) {
	m, err := decode[M](mr.collection.FindOne(ctx, bson.M{
		"_id": id}))
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = liara.ErrNotFound
	}

	v, e := mr.mapper.FromM(&m)
	return v, e, err
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
	f any,
	s any,
	p Page,
) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		result, err := mr.collection.Find(ctx, f, options.Find().
			SetSort(s).
			SetSkip(int64(p.Offset)).
			SetLimit(int64(p.Limit)))
		if err != nil {
			yield(nil, err)
			return
		}

		defer result.Close(ctx)

		for result.Next(ctx) {
			m, err := decode[M](result)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					err = liara.ErrNotFound
				}
				yield(nil, err)
				return
			}

			_, e := mr.mapper.FromM(&m)
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
