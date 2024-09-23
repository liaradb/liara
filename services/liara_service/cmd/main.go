package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/cardboardrobots/essql"
	"github.com/cardboardrobots/eventsource"
	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara_service/config"
	"github.com/cardboardrobots/liara_service/feature/base"
	"github.com/cardboardrobots/liara_service/feature/eventsource/controller"
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

	_ = eventsource.MockEventSource{}
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
	eventRepository := essql.NewEventRepository(db, "events")
	err := eventRepository.CreateTable(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = eventRepository.CreateIndex(ctx)
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterEventSourceServiceServer(s, controller.NewEventSourceController(
		service.NewEventService(eventRepository),
	))

	return s
}

func ConnectPostgresDB(uri string) (*sql.DB, error) {
	return sql.Open("postgres", uri)
}

func ConnectSqliteDB(uri string) (*sql.DB, error) {
	return sql.Open("sqlite", uri)
}
