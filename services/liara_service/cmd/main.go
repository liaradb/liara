package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	esservice "github.com/cardboardrobots/eventsource/service"
	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara_service/config"
	"github.com/cardboardrobots/liara_service/feature/base"
	"github.com/cardboardrobots/liara_service/feature/eventsource/controller"
	"github.com/cardboardrobots/liara_service/feature/eventsource/infrastructure"
	"github.com/cardboardrobots/liara_service/feature/eventsource/service"
	"github.com/cardboardrobots/listener"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

func main() {
	log.SetPrefix("[liara]\t")
	log.Println("started...")

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := ConnectSqliteDB(conf.SqliteDbUri)
	if err != nil {
		log.Fatal(err)
	}

	listener.Listen(context.Background(), conf.Port, conf.Port+1,
		http.NewServeMux(),
		initService(db))

	_ = esservice.MockEventSource{}
}

func initService(db *sql.DB) *grpc.Server {
	s := listener.NewServerBuilder().
		AddUnary(
			listener.LogGRPC(false),
			listener.ErrorInterceptor(base.GetStatusCodeGRPC),
		).
		AddStream(
			listener.LogStreamGRPC(false),
			listener.ErrorInterceptorStream(base.GetStatusCodeGRPC),
		).
		Build()

	ctx := context.Background()
	eventRepository, outboxRepository, err := createTable(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterEventSourceServiceServer(s, controller.NewEventSourceController(
		service.NewEventService(eventRepository, eventRepository, outboxRepository),
	))

	return s
}

func createTable(ctx context.Context, db *sql.DB) (*infrastructure.EventRepository, *infrastructure.OutboxRepository, error) {
	eventRepository := infrastructure.NewEventRepository(db, "events")
	err := eventRepository.CreateTable(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = eventRepository.CreateIndex(ctx)
	if err != nil {
		return nil, nil, err
	}

	outboxRepository := infrastructure.NewOutboxRepository(db, "outbox")
	err = outboxRepository.CreateTable(ctx)
	if err != nil {
		return nil, nil, err
	}

	return &eventRepository, outboxRepository, nil
}

func ConnectPostgresDB(uri string) (*sql.DB, error) {
	return sql.Open("postgres", uri)
}

func ConnectSqliteDB(uri string) (*sql.DB, error) {
	return sql.Open("sqlite", uri)
}
