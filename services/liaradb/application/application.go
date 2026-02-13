package application

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"time"

	"github.com/cardboardrobots/errormap"
	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/application/listener"
	"github.com/liaradb/liaradb/collection"
	"github.com/liaradb/liaradb/collection/btree"
	"github.com/liaradb/liaradb/collection/tablename"
	"github.com/liaradb/liaradb/collection/tenant"
	"github.com/liaradb/liaradb/controller"
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/service"
	"github.com/liaradb/liaradb/encoder/raw"
	"github.com/liaradb/liaradb/file/disk"
	"github.com/liaradb/liaradb/locktable"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/transaction"
	"google.golang.org/grpc"
)

type Application struct {
	conf        configuration
	storage     *storage.Storage
	collections *collection.Collections
	txManager   *transaction.Manager
	log         *recovery.Log
	lockTable   *locktable.LockTable[action.ItemID] // TODO: Is this ID type correct?
}

func New(conf configuration) *Application {
	segmentSize := 1024
	inSize := 100

	fsys := &disk.FileSystem{}

	s := storage.New(fsys, conf.Buffers, int64(conf.BlockSize), path.Join(conf.Directory, "table"))
	log := recovery.NewLog(int64(conf.BlockSize), action.PageID(segmentSize), fsys, path.Join(conf.Directory, "log"))
	lt := locktable.NewLockTable[action.ItemID](inSize)

	return &Application{
		conf:        conf,
		storage:     s,
		collections: collection.NewCollections(s),
		txManager:   transaction.NewManager(log, s, lt),
		log:         log,
		lockTable:   lt,
	}
}

// TODO: Ensure all goroutines are stopped before calling close
func (a *Application) Run(ctx context.Context) error {
	ctx, cancelMain := context.WithCancel(ctx)
	if err := a.run(ctx); err != nil {
		cancelMain()
		return err
	}

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

	a.txManager.Run(ctx)

	if err := a.storage.Run(ctx); err != nil {
		return err
	}

	slog.Info("storage running")

	if err := a.log.Open(ctx); err != nil {
		return err
	}

	slog.Info("recovering...")

	if err := a.recover(ctx); err != nil {
		return err
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

func (a *Application) recover(ctx context.Context) error {
	it, err := a.log.Recover()
	if err != nil {
		return err
	}

	for r := range it {
		switch r.Action() {
		case record.ActionCheckpoint:
			fmt.Printf("recover: %v\n", r.Action())
		case record.ActionCommit:
			fmt.Printf("recover: %v\n", r.Action())
		case record.ActionInsert:
			if err := a.recoverEvent(ctx, r); err != nil {
				return err
			}
		case record.ActionRemove:
			fmt.Printf("recover: %v\n", r.Action())
		case record.ActionRollback:
			fmt.Printf("recover: %v\n", r.Action())
		case record.ActionUpdate:
			fmt.Printf("recover: %v\n", r.Action())
		default:
			fmt.Printf("recover: %v\n", r.Action())
		}
	}

	return nil
}

func (a *Application) recoverEvent(ctx context.Context, r *record.Record) error {
	var e entity.Event
	if err := e.Read(raw.NewBufferFromSlice(r.Data())); err != nil {
		return err
	}

	fmt.Printf("recover: %v: %v\n", r.Action(), e.AggregateID.String())
	tn := tablename.New(r.TenantID())
	_, err := a.collections.EventLog.Append(ctx, tn, e.PartitionID, &e)
	return err
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
	if err := a.txManager.Flush(time.Now()); err != nil {
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

	_, err := a.createRepositories()
	if err != nil {
		slog.Error("create repositories",
			"error", err)
		panic(err)
	}

	pb.RegisterEventSourceServiceServer(s, controller.NewEventSourceController(
		service.NewEventService(
			a.txManager,
		),
		service.NewTenantService(
			tenant.New(a.storage, btree.NewCursor(a.storage))),
	))

	return s
}

type repositories struct {
}

func (a *Application) createRepositories() (*repositories, error) {
	return &repositories{}, nil
}
