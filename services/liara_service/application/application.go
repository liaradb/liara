package application

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/cardboardrobots/errormap"
	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara_service/feature/eventsource/controller"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/service"
	"github.com/cardboardrobots/liara_service/feature/eventsource/infrastructure"
	"github.com/cardboardrobots/listener"
	"google.golang.org/grpc"
)

func Run() error {
	conf, err := LoadConfig()
	if err != nil {
		return err
	}

	db, err := ConnectSqliteDB(conf.SqliteDbUri)
	if err != nil {
		return err
	}

	listener.Listen(context.Background(), conf.Port, conf.Port+1,
		http.NewServeMux(),
		initService(db))

	return nil
}

func initService(db *sql.DB) *grpc.Server {
	s := listener.NewServerBuilder().
		AddUnary(
			listener.LogGRPC(false),
			listener.ErrorInterceptor(errormap.GetStatusCodeGRPC),
		).
		AddStream(
			listener.LogStreamGRPC(false),
			listener.ErrorInterceptorStream(errormap.GetStatusCodeGRPC),
		).
		Build()

	ctx := context.Background()
	r, err := createRepositories(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	pb.RegisterEventSourceServiceServer(s, controller.NewEventSourceController(
		service.NewEventService(
			r.transactionContainer,
			r.eventRepository,
			r.outboxRepository,
			r.requestRepository),
		service.NewTenantService(
			r.transactionContainer,
			r.eventRepository,
			r.outboxRepository,
			r.requestRepository,
			r.tenantRepository),
	))

	return s
}

type repositories struct {
	transactionContainer *infrastructure.TransactionContainer
	eventRepository      *infrastructure.EventRepository
	outboxRepository     *infrastructure.OutboxRepository
	requestRepository    *infrastructure.RequestRepository
	tenantRepository     *infrastructure.TenantRepository
}

func createRepositories(ctx context.Context, db *sql.DB) (*repositories, error) {
	transactionContainer := infrastructure.NewTransactionContainer(db, &sql.TxOptions{Isolation: sql.LevelSerializable})

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
		transactionContainer: transactionContainer,
		eventRepository:      eventRepository,
		outboxRepository:     outboxRepository,
		requestRepository:    requestRepository,
		tenantRepository:     tenantRepository,
	}, nil
}

func ConnectPostgresDB(uri string) (*sql.DB, error) {
	return sql.Open("postgres", uri)
}

func ConnectSqliteDB(uri string) (*sql.DB, error) {
	return sql.Open("sqlite", uri)
}
