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

func Insert[T any, I ~string, M any](
	ctx context.Context,
	c *mongo.Collection,
	toM func(liara.Version, I, *T) *M,
	v liara.Version,
	id I,
	t *T,
) error {
	_, err := c.ReplaceOne(ctx,
		bson.M{"_id": id},
		toM(v, id, t),
		options.Replace().SetUpsert(true))
	return err
}

func Get[T any, I ~string, M any](
	ctx context.Context,
	c *mongo.Collection,
	fromM func(*M) (liara.Version, *T),
	id I,
) (liara.Version, *T, error) {
	var m M
	err := c.FindOne(ctx,
		bson.M{"_id": id}).
		Decode(&m)
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = liara.ErrNotFound
	}

	v, e := fromM(&m)
	return v, e, err
}

func Delete[I ~string](
	ctx context.Context,
	c *mongo.Collection,
	id I,
) error {
	_, err := c.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func GetList[T any, F any, S any, M any](
	ctx context.Context,
	c *mongo.Collection,
	fromM func(*M) (liara.Version, *T),
	f F,
	s S,
	p Page,
) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		result, err := c.Find(ctx, f, options.Find().
			SetSort(s).
			SetSkip(int64(p.Offset)).
			SetLimit(int64(p.Limit)))
		if err != nil {
			yield(nil, err)
			return
		}

		defer result.Close(ctx)

		for result.Next(ctx) {
			var m M
			err = result.Decode(&m)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					err = liara.ErrNotFound
				}
				yield(nil, err)
				return
			}

			_, e := fromM(&m)
			if yield(e, nil) {
				return
			}
		}
	}
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

func Watch[P []bson.D](
	ctx context.Context,
	coll *mongo.Collection,
	pipeline P,
	token string,
) (iter.Seq2[bson.Raw, string], error) {
	o := options.ChangeStream()
	if token != "" {
		o.SetResumeAfter(bson.Raw(token))
	}

	cs, err := coll.Watch(ctx, pipeline, o)
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
