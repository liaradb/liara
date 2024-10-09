package esmongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/cardboardrobots/eventsource"
	"github.com/cardboardrobots/eventsource/value"
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
	toM func(value.Version, I, *T) *M,
	v value.Version,
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
	fromM func(*M) (value.Version, *T),
	id I,
) (value.Version, *T, error) {
	var m M
	err := c.FindOne(ctx,
		bson.M{"_id": id}).
		Decode(&m)
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = eventsource.ErrNotFound
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
	fromM func(*M) (value.Version, *T),
	f F,
	s S,
	p Page,
	callback func(*T) error,
) (int, error) {
	count, err := c.CountDocuments(ctx, f)
	if err != nil {
		return 0, err
	}

	result, err := c.Find(ctx, f, options.Find().
		SetSort(s).
		SetSkip(int64(p.Offset)).
		SetLimit(int64(p.Limit)))
	if err != nil {
		return int(count), err
	}

	defer result.Close(ctx)

	for result.Next(ctx) {
		var m M
		err = result.Decode(&m)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				err = eventsource.ErrNotFound
			}
			return int(count), err
		}

		_, e := fromM(&m)
		err = callback(e)
		if err != nil {
			return int(count), err
		}
	}

	return int(count), nil
}

func GetListSlice[T any, F any, S any, M any](
	ctx context.Context,
	c *mongo.Collection,
	fromM func(*M) (value.Version, *T),
	f F,
	s S,
	p Page,
) ([]*T, int, error) {
	count, err := c.CountDocuments(ctx, f)
	if err != nil {
		return nil, 0, err
	}

	m := []M{}
	result, err := c.Find(ctx, f, options.Find().
		SetSort(s).
		SetSkip(int64(p.Offset)).
		SetLimit(int64(p.Limit)))
	if err != nil {
		return nil, 0, err
	}

	err = result.All(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	t := make([]*T, 0, len(m))
	for _, m := range m {
		_, e := fromM(&m)
		t = append(t, e)
	}

	return t, int(count), nil
}

func RunTransaction[T any](ctx context.Context, c *mongo.Client, p func(ctx context.Context) (T, error)) (T, error) {
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

func Watch(ctx context.Context, coll *mongo.Collection, token string) (string, error) {
	o := options.ChangeStream()
	if token != "" {
		resumeToken := bson.Raw(token)
		o.SetResumeAfter(resumeToken)
	}

	cs, err := coll.Watch(ctx, mongo.Pipeline{}, o)
	if err != nil {
		return token, err
	}

	defer cs.Close(ctx)

	for cs.Next(ctx) {
		next := cs.Current
		fmt.Println(next)

		token = cs.ResumeToken().String()
		fmt.Println(token)
	}

	return token, err
}
