package esmongo

import (
	"context"
	"iter"

	"go.mongodb.org/mongo-driver/mongo"
)

type Page struct {
	Offset int
	Limit  int
}

type Mapper[T any, M any] interface {
	FromModel(*M) *T
	ToModel(*T) *M
}

type Repository[T any, I ~string, M any] struct {
	collection *Collection[Model[M]]
	mapper     Mapper[T, Model[M]]
}

func NewRepository[T any, I ~string, M any](
	collection *mongo.Collection,
	mapper Mapper[T, Model[M]],
) *Repository[T, I, M] {
	return &Repository[T, I, M]{
		collection: NewCollection[Model[M]](collection),
		mapper:     mapper,
	}
}

func (r *Repository[T, I, M]) Insert(ctx context.Context, id I, t *T) error {
	return r.collection.Insert(ctx,
		string(id),
		r.mapper.ToModel(t))
}

func (r *Repository[T, I, M]) Replace(ctx context.Context, id I, v int, t *T) error {
	return r.collection.Replace(ctx,
		Filter().
			Property("_id", id).
			Property("version", v),
		r.mapper.ToModel(t))
}

func (r *Repository[T, I, M]) Get(ctx context.Context, id I) (*T, error) {
	m, err := r.collection.Get(ctx, string(id))
	return r.mapper.FromModel(m), err
}

func (r *Repository[T, I, M]) Delete(ctx context.Context, id I) error {
	return r.collection.Delete(ctx, string(id))
}

func (r *Repository[T, I, M]) GetList(ctx context.Context, filter FilterBuilder, sort *SortBuilder) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		for row, err := range r.collection.GetList(ctx, filter, sort) {
			if !yield(r.mapper.FromModel(row), err) {
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
