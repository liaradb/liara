package application

import (
	"context"
	"log/slog"
	"path"

	"github.com/cardboardrobots/errormap"
	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/application/listener"
	"github.com/liaradb/liaradb/collection/eventlog"
	"github.com/liaradb/liaradb/controller"
	"github.com/liaradb/liaradb/domain/infrastructure"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/file/disk"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/page"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction"
	"google.golang.org/grpc"
)

type Application struct {
	conf      configuration
	eventLog  *eventlog.EventLog
	storage   *storage.Storage
	txManager *transaction.Manager
	log       *recovery.Log
	lockTable *locktable.LockTable[action.ItemID] // TODO: Is this ID type correct?
}

func New(conf configuration) *Application {
	segmentSize := 1024
	inSize := 100

	fsys := &disk.FileSystem{}

	s := storage.New(fsys, conf.Buffers, int64(conf.BlockSize), path.Join(conf.Directory, "table"))
	log := recovery.NewLog(int64(conf.BlockSize), page.PageID(segmentSize), fsys, path.Join(conf.Directory, "log"))
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

// TODO: Ensure all goroutines are stopped before calling close
func (a *Application) Run(ctx context.Context) error {
	ctx, cancelMain := context.WithCancel(ctx)
	a.run(ctx)
	defer a.close()
	defer func() {
		slog.Info("shutting down...")
		cancelMain()
	}()

	ctx, cancelListen := WithSignal(ctx)
	defer cancelListen()

	a.listen(ctx)

	return nil
}

func (a *Application) run(ctx context.Context) error {
	slog.Info("starting...")

	if err := a.storage.Run(ctx); err != nil {
		return err
	}

	slog.Info("storage running")

	if err := a.log.Open(ctx); err != nil {
		return err
	}

	slog.Info("recovering...")

	it, err := a.log.Recover()
	if err != nil {
		return err
	}

	for range it {
		// fmt.Printf("recover: %v\n", r.Action())
	}

	slog.Info("recovered")

	if err := a.log.StartWriter(); err != nil {
		return err
	}

	slog.Info("log running")

	a.lockTable.Run(ctx)

	slog.Info("lock table running")

	return nil
}

func (a *Application) listen(ctx context.Context) {
	listener.Listen(ctx, a.conf.Port, a.initService())
}

// Closing Process
//   - close gRPC requests
//   - Cancel Context
//   - Flush Log
//   - Flush Buffers
func (a *Application) close() {
	slog.Info("flushing...")
	if err := a.storage.FlushAll(); err != nil {
		slog.Error("unable to flush",
			"error", err)
		return
	}
	slog.Info("flushing complete")

	slog.Info("shutdown complete")
}

func (a *Application) initService() *grpc.Server {
	s := listener.NewServerBuilder().
		AddUnary(
			listener.LogGRPC(),
			listener.ErrorInterceptor(errormap.GetStatusCodeGRPC),
		).
		AddStream(
			listener.LogStreamGRPC(),
			listener.ErrorInterceptorStream(errormap.GetStatusCodeGRPC),
		).
		Build()

	r, err := a.createRepositories()
	if err != nil {
		slog.Error("create repositories",
			"error", err)
		panic(err)
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
			nil,
			a.eventLog,
			nil,
			"testfile"),
	}, nil
}
