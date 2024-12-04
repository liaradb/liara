package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara_service/config"
	"github.com/cardboardrobots/liara_service/feature/base"
	"github.com/cardboardrobots/liara_service/feature/eventsource/controller"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/service"
	"github.com/cardboardrobots/liara_service/feature/eventsource/infrastructure"
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
	r, err := createRepositories(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterEventSourceServiceServer(s, controller.NewEventSourceController(
		service.NewEventService(
			r.transactionRepository,
			r.eventRepository,
			r.outboxRepository,
			r.requestRepository),
		service.NewTenantService(
			r.transactionRepository,
			r.eventRepository,
			r.outboxRepository,
			r.requestRepository,
			r.tenantRepository),
	))

	return s
}

type repositories struct {
	transactionRepository *infrastructure.TransactionRepository
	eventRepository       *infrastructure.EventRepository
	outboxRepository      *infrastructure.OutboxRepository
	requestRepository     *infrastructure.RequestRepository
	tenantRepository      *infrastructure.TenantRepository
}

func createRepositories(ctx context.Context, db *sql.DB) (*repositories, error) {
	transactionRepository := infrastructure.NewTransactionRepository(db, &sql.TxOptions{Isolation: sql.LevelSerializable})

	eventRepository := infrastructure.NewEventRepository(db)
	err := eventRepository.CreateTable(ctx, "")
	if err != nil {
		return nil, err
	}

	err = eventRepository.CreateIndex(ctx, "")
	if err != nil {
		return nil, err
	}

	outboxRepository := infrastructure.NewOutboxRepository(db)
	err = outboxRepository.CreateTable(ctx, "")
	if err != nil {
		return nil, err
	}

	requestRepository := infrastructure.NewRequestRepository(db)
	err = requestRepository.CreateTable(ctx, "")
	if err != nil {
		return nil, err
	}

	tenantRepository := infrastructure.NewTenantRepository(db)
	err = tenantRepository.CreateTable(ctx)
	if err != nil {
		return nil, err
	}

	return &repositories{
		transactionRepository: transactionRepository,
		eventRepository:       &eventRepository,
		outboxRepository:      outboxRepository,
		requestRepository:     requestRepository,
		tenantRepository:      tenantRepository,
	}, nil
}

func ConnectPostgresDB(uri string) (*sql.DB, error) {
	return sql.Open("postgres", uri)
}

func ConnectSqliteDB(uri string) (*sql.DB, error) {
	return sql.Open("sqlite", uri)
}
