package application

import (
	"context"
	"log"
	"net/http"

	"github.com/cardboardrobots/errormap"
	"github.com/cardboardrobots/listener"
	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/controller"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/file/disk"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/eventlog"
	"google.golang.org/grpc"
)

type Application struct {
	eventLog *eventlog.EventLog
	storage  *storage.Storage
}

func New(max int, bs int64) *Application {
	fsys := &disk.FileSystem{}
	storage := storage.NewStorage(fsys, max, bs)
	return &Application{
		eventLog: eventlog.New(storage),
		storage:  storage,
	}
}

func (a *Application) Run(ctx context.Context) error {
	conf, err := LoadConfig()
	if err != nil {
		return err
	}

	a.storage.Run(ctx)

	listener.Listen(ctx, conf.Port, conf.Port+1,
		http.NewServeMux(),
		initService())

	return nil
}

func initService() *grpc.Server {
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

	r, err := createRepositories()
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
	transactionContainer service.TransactionContainer
	eventRepository      service.EventRepository
	outboxRepository     service.OutboxRepository
	requestRepository    service.RequestRepository
	tenantRepository     service.TenantRepository
}

func createRepositories() (*repositories, error) {
	return &repositories{}, nil
}
