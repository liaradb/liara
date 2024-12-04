package infrastructure

import (
	"context"
	"reflect"
	"testing"

	"database/sql"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/service"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
	_ "modernc.org/sqlite"
)

func connectSqliteDB(uri string) (*sql.DB, error) {
	return sql.Open("sqlite", uri)
}

func connectInMemory(ctx context.Context, tenantID value.TenantID) (*EventRepository, error) {
	db, err := connectSqliteDB(":memory:")
	if err != nil {
		return nil, err
	}

	er := NewEventRepository(db)

	if err = er.CreateTable(ctx, tenantID); err != nil {
		return nil, err
	}

	if err = er.CreateIndex(ctx, tenantID); err != nil {
		return nil, err
	}

	return &er, nil
}

func TestEventRepository_Append(t *testing.T) {
	t.Run("should append event", func(t *testing.T) {
		ctx := context.Background()

		er, err := connectInMemory(ctx, "")
		if err != nil {
			t.Fatal(err)
		}

		aggregateID := value.AggregateID("aggregateID")

		want := entity.Event{
			GlobalVersion: 1,
			AggregateName: "example",
			ID:            "eventID",
			AggregateID:   aggregateID,
			Version:       1,
		}

		err = er.Append(ctx, "", service.AppendEvent{
			AggregateName: "example",
			ID:            "eventID",
			AggregateID:   aggregateID,
			Version:       1,
		})
		if err != nil {
			t.Fatal(err)
		}

		rows := er.Get(ctx, "", aggregateID)
		count := 0
		for e, err := range rows {
			if err != nil {
				t.Error(err)
			}

			count++
			if !reflect.DeepEqual(e, want) {
				t.Error("event is incorrect")
			}
		}

		if count != 1 {
			t.Errorf("count is incorrect.  Recieved: %v, Expected: %v", count, 1)
		}
	})

	t.Run("should not append existing version", func(t *testing.T) {
		ctx := context.Background()

		er, err := connectInMemory(ctx, "")
		if err != nil {
			t.Fatal(err)
		}

		aggregateID := value.AggregateID("aggregateID")

		want := service.AppendEvent{
			AggregateName: "example",
			ID:            "eventID",
			AggregateID:   aggregateID,
			Version:       1,
		}

		err = er.Append(ctx, "", want)
		if err != nil {
			t.Fatal(err)
		}

		err = er.Append(ctx, "", want)
		if err == nil {
			t.Error("should return error")
		}
	})
}
