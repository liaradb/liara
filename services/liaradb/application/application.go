package application

import (
	"context"
	l "log"
	"net/http"

	"github.com/cardboardrobots/errormap"
	"github.com/cardboardrobots/listener"
	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/controller"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/file/disk"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/log"
	"github.com/liaradb/liaradb/log/action"
	"github.com/liaradb/liaradb/log/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/eventlog"
	"github.com/liaradb/liaradb/transaction"
	"google.golang.org/grpc"
)

type Application struct {
	eventLog  *eventlog.EventLog
	storage   *storage.Storage
	txManager *transaction.Manager
	log       *log.Log
	lockTable *locktable.LockTable[action.ItemID] // TODO: Is this ID type correct?
}

func New(max int, bs int64) *Application {
	segmentSize := 1024
	inSize := 100

	fsys := &disk.FileSystem{}

	s := storage.NewStorage(fsys, max, bs, ".dbdata/table")
	log := log.NewLog(bs, page.PageID(segmentSize), fsys, ".dbdata/log")
	lt := locktable.NewLockTable[action.ItemID](inSize)

	return &Application{
		eventLog:  eventlog.New(s),
		storage:   s,
		txManager: transaction.NewManager(log, s, lt),
		log:       log,
		lockTable: lt,
	}
}

func (a *Application) Run(ctx context.Context) error {
	conf, err := LoadConfig()
	if err != nil {
		return err
	}

	if err := a.storage.Run(ctx); err != nil {
		return err
	}

	if err := a.log.Open(ctx); err != nil {
		return err
	}

	a.lockTable.Run(ctx)

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
		l.Fatal(err)
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
