package application

import (
	"context"
	l "log"
	"net/http"
	"path"

	"github.com/cardboardrobots/errormap"
	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/controller"
	"github.com/liaradb/liaradb/domain/infrastructure"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/file/disk"
	"github.com/liaradb/liaradb/listener"
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
	conf      configuration
	eventLog  *eventlog.EventLog
	storage   *storage.Storage
	txManager *transaction.Manager
	log       *log.Log
	lockTable *locktable.LockTable[action.ItemID] // TODO: Is this ID type correct?
}

func New(conf configuration) *Application {
	segmentSize := 1024
	inSize := 100

	fsys := &disk.FileSystem{}

	s := storage.NewStorage(fsys, conf.Buffers, int64(conf.BlockSize), path.Join(conf.Directory, "table"))
	log := log.NewLog(int64(conf.BlockSize), page.PageID(segmentSize), fsys, path.Join(conf.Directory, "log"))
	lt := locktable.NewLockTable[action.ItemID](inSize)

	return &Application{
		conf:      conf,
		eventLog:  eventlog.New(s),
		storage:   s,
		txManager: transaction.NewManager(log, s, lt),
		log:       log,
		lockTable: lt,
	}
}

func (a *Application) Run(ctx context.Context) error {
	if err := a.storage.Run(ctx); err != nil {
		return err
	}

	defer a.Close()

	if err := a.log.Open(ctx); err != nil {
		return err
	}

	if err := a.log.StartWriter(); err != nil {
		return err
	}

	a.lockTable.Run(ctx)

	listener.Listen(ctx, a.conf.Port, a.conf.Port+1,
		http.NewServeMux(),
		a.initService())

	return nil
}

func (a *Application) Close() {
	l.Println("shutting down...")

	l.Println("flushing...")
	if err := a.storage.FlushAll(); err != nil {
		l.Fatal(err)
	}
	l.Println("flushing complete")
}

func (a *Application) initService() *grpc.Server {
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

	r, err := a.createRepositories()
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

func (a *Application) createRepositories() (*repositories, error) {
	return &repositories{
		// TODO: Change the file name
		eventRepository: infrastructure.NewEventRepository(
			a.txManager,
			a.eventLog,
			"testfile"),
	}, nil
}
