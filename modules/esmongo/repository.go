package esmongo

import (
	"context"
	"iter"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository[I EntityID, E Entity[I], M any, V Event] struct {
	collection *Collection[model[M]]
	mapper     Mapper[I, E, M, V]
}

type Mapper[I EntityID, E Entity[I], M any, V Event] interface {
	FromModel(string, int, M) E
	ToModel(E) M
	FromModelEvent(V) ([]byte, error)
	ToModelEvent(string, []byte) (V, bool)
}

func NewRepository[I EntityID, E Entity[I], M any, V Event](
	collection *mongo.Collection,
	mapper Mapper[I, E, M, V],
) *Repository[I, E, M, V] {
	return &Repository[I, E, M, V]{
		collection: NewCollection[model[M]](collection),
		mapper:     mapper,
	}
}

func (r *Repository[I, E, M, V]) Insert(
	ctx context.Context,
	entity E,
	events []V,
) error {
	return r.collection.Insert(ctx,
		entity.ID().String(),
		r.newModel(entity, events))
}

func (r *Repository[I, E, M, V]) Replace(
	ctx context.Context,
	entity E,
	events []V,
) error {
	return r.collection.Replace(ctx,
		Filter().
			Property("_id", entity.ID().String()).
			Property("version", entity.Version()),
		r.newModel(entity, events).
			increment())
}

func (r *Repository[I, E, M, V]) ReplaceAtVersion(
	ctx context.Context,
	version int,
	entity E,
	events []V,
) error {
	return r.collection.Replace(ctx,
		Filter().
			Property("_id", entity.ID().String()).
			Property("version", version),
		r.newModel(entity, events))
}

func (r *Repository[I, E, M, V]) newModel(entity E, events []V) *model[M] {
	m := r.mapper.ToModel(entity)
	evs := make([]*modelEvent, 0, len(events))
	for _, ev := range events {
		e, _ := r.mapper.FromModelEvent(ev)
		t := ev.Type()
		evs = append(evs, newModelEvent(t, e))
	}
	return newModel(
		entity.ID().String(),
		entity.Version(),
		m,
		evs)
}

func (r *Repository[I, E, M, V]) Get(
	ctx context.Context,
	id I,
) (E, error) {
	m, err := r.collection.Get(ctx, id.String())
	return r.mapper.FromModel(m.ID, m.Version, m.Value), err
}

func (r *Repository[I, E, M, V]) Delete(
	ctx context.Context,
	id I,
) error {
	return r.collection.Delete(ctx, id.String())
}

func (r *Repository[I, E, M, V]) GetList(
	ctx context.Context,
	filter FilterBuilder,
	sort *SortBuilder,
) iter.Seq2[E, error] {
	return func(yield func(E, error) bool) {
		for row, err := range r.collection.GetList(ctx, filter, sort) {
			if !yield(r.mapper.FromModel(row.ID, row.Version, row.Value), err) {
				return
			}
		}
	}
}

func (r *Repository[I, E, M, V]) Watch(ctx context.Context, pipeline any, token string) iter.Seq2[Change[[]any], error] {
	return func(yield func(Change[[]any], error) bool) {
		rows := r.collection.Watch(ctx, pipeline, token)

		for row, err := range rows {
			if err != nil {
				yield(Change[[]any]{Token: row.Token}, err)
				return
			}

			events := make([]any, 0, len(row.Value.Events))
			for _, res := range row.Value.Events {
				if e, ok := r.mapper.ToModelEvent(res.Type, res.Data); ok {
					events = append(events, e)
				}
			}
			if !yield(Change[[]any]{Value: events, Token: row.Token}, nil) {
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
