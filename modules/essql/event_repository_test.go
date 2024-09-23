package essql

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"database/sql"

	"github.com/cardboardrobots/eventsource"
	_ "modernc.org/sqlite"
)

func connectSqliteDB(uri string) (*sql.DB, error) {
	return sql.Open("sqlite", uri)
}

func connectInMemory(ctx context.Context, name string) (*EventRepository, error) {
	db, err := connectSqliteDB(":memory:")
	if err != nil {
		return nil, err
	}

	er := NewEventRepository(db, name)

	if err = er.CreateTable(ctx); err != nil {
		return nil, err
	}

	if err = er.CreateIndex(ctx); err != nil {
		return nil, err
	}

	return &er, nil
}

func TestEventRepository_Append(t *testing.T) {
	t.Run("should append event", func(t *testing.T) {
		ctx := context.Background()

		er, err := connectInMemory(ctx, "events")
		if err != nil {
			t.Fatal(err)
		}

		aggregateID := eventsource.AggregateID("aggregateID")

		want := eventsource.Event{
			GlobalVersion: 1,
			AggregateName: "example",
			ID:            "eventID",
			AggregateID:   aggregateID,
			Version:       1,
		}

		err = er.Append(ctx, want)
		if err != nil {
			t.Fatal(err)
		}

		rows := er.Get(ctx, aggregateID)
		count := 0
		for e, err := range rows {
			if err != nil {
				t.Error(err)
			}

			count++
			if reflect.DeepEqual(e, want) {
				t.Error("event is incorrect")
			}
		}

		if count != 1 {
			t.Errorf("count is incorrect.  Recieved: %v, Expected: %v", count, 1)
		}
	})

	t.Run("should not append invalid version", func(t *testing.T) {
		ctx := context.Background()

		er, err := connectInMemory(ctx, "events")
		if err != nil {
			t.Fatal(err)
		}

		aggregateID := eventsource.AggregateID("aggregateID")

		want := eventsource.Event{
			GlobalVersion: 1,
			AggregateName: "example",
			ID:            "eventID",
			AggregateID:   aggregateID,
			Version:       0,
		}

		err = er.Append(ctx, want)
		if !errors.Is(err, eventsource.ErrAggregateVersionInvalid) {
			t.Error("should return error")
		}
	})

	t.Run("should not append existing version", func(t *testing.T) {
		ctx := context.Background()

		er, err := connectInMemory(ctx, "events")
		if err != nil {
			t.Fatal(err)
		}

		aggregateID := eventsource.AggregateID("aggregateID")

		want := eventsource.Event{
			GlobalVersion: 1,
			AggregateName: "example",
			ID:            "eventID",
			AggregateID:   aggregateID,
			Version:       1,
		}

		err = er.Append(ctx, want)
		if err != nil {
			t.Fatal(err)
		}

		err = er.Append(ctx, want)
		if err == nil {
			t.Error("should return error")
		}
	})
}
