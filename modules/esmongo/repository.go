package esmongo

import (
	"context"
	"iter"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository[I EntityID, E Entity[I], M any] struct {
	collection *Collection[model[M]]
	mapper     Mapper[I, E, M]
}

type model[T any] struct {
	Record `bson:"inline"`
	Value  T `bson:"inline"`
}

type Mapper[I EntityID, E Entity[I], M any] interface {
	FromModel(Record, M) E
	ToModel(E) M
	ToRecordEvent(Event) (RecordEvent, bool)
	FromRecordEvent(RecordEvent) (Event, bool)
}

func NewRepository[I EntityID, E Entity[I], M any](
	collection *mongo.Collection,
	mapper Mapper[I, E, M],
) *Repository[I, E, M] {
	return &Repository[I, E, M]{
		collection: NewCollection[model[M]](collection),
		mapper:     mapper,
	}
}

func (r *Repository[I, E, M]) Insert(ctx context.Context, e E) error {
	return r.collection.Insert(ctx,
		e.ID().String(),
		r.newModel(e))
}

func (r *Repository[I, E, M]) Replace(ctx context.Context, id I, v int, e E) error {
	return r.collection.Replace(ctx,
		Filter().
			Property("_id", id.String()).
			Property("version", v),
		r.newModel(e))
}

func (r *Repository[I, E, M]) newModel(entity E) *model[M] {
	m := r.mapper.ToModel(entity)
	events, _ := newRecordEvents(entity.Events())
	return &model[M]{
		Record: Record{
			ID:      entity.ID().String(),
			Version: entity.Version().Value(),
			Events:  events,
		},
		Value: m,
	}
}

func (r *Repository[I, E, M]) Get(ctx context.Context, id I) (E, error) {
	m, err := r.collection.Get(ctx, id.String())
	return r.mapper.FromModel(m.Record, m.Value), err
}

func (r *Repository[I, E, M]) Delete(ctx context.Context, id I) error {
	return r.collection.Delete(ctx, id.String())
}

func (r *Repository[I, E, M]) GetList(ctx context.Context, filter FilterBuilder, sort *SortBuilder) iter.Seq2[E, error] {
	return func(yield func(E, error) bool) {
		for row, err := range r.collection.GetList(ctx, filter, sort) {
			if !yield(r.mapper.FromModel(row.Record, row.Value), err) {
				return
			}
		}
	}
}

func (r *Repository[I, E, M]) Watch(ctx context.Context, pipeline any, token string) iter.Seq2[Change[[]Event], error] {
	return func(yield func(Change[[]Event], error) bool) {
		rows := r.collection.Watch(ctx, pipeline, token)

		for row, err := range rows {
			if err != nil {
				yield(Change[[]Event]{Token: row.Token}, err)
				return
			}

			events := make([]Event, 0, len(row.Value.Events))
			for _, res := range row.Value.Events {
				if e, ok := r.mapper.FromRecordEvent(*res); ok {
					events = append(events, e)
				}
			}
			if !yield(Change[[]Event]{Value: events, Token: row.Token}, nil) {
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
