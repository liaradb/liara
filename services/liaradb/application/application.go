package application

import (
	"context"
	"log/slog"
	"path"
	"time"

	"github.com/cardboardrobots/errormap"
	pb "github.com/liaradb/eventsource_go/generated"
	"github.com/liaradb/liaradb/application/listener"
	"github.com/liaradb/liaradb/controller"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/lrupool"
	"github.com/liaradb/liaradb/transaction"
	"github.com/liaradb/liaradb/transaction/log"
	"google.golang.org/grpc"
)

type Application struct {
	conf      configuration
	fsys      *filecache.Cache
	s         *grpc.Server
	txManager *transaction.Manager
	services  *services
	ot        openTelemetry
}

const (
	segmentSize = 1024
	inSize      = 100
)

func New(conf configuration) *Application {
	fsys := filecache.New()

	s := storage.New(
		fsys,
		lrupool.New(),
		conf.Buffers,
		int16(conf.BlockSize),
		path.Join(conf.Directory, "tenants"))

	log := log.New(
		int16(conf.BlockSize),
		segment.PageID(segmentSize),
		int64(conf.RecordSize),
		conf.WriteQueueSize,
		fsys,
		path.Join(conf.Directory, "log"))

	txManager := transaction.NewManager(log, s, inSize)

	return &Application{
		conf:      conf,
		fsys:      fsys,
		txManager: txManager,
		services:  newServices(txManager, s),
	}
}

// TODO: Ensure all goroutines are stopped before calling close
func (a *Application) Run(ctx context.Context) error {
	ctxMain, cancelMain := context.WithCancel(ctx)

	a.ot.init(ctxMain)

	a.txManager.Run(ctxMain)
	if err := a.txManager.Recover(ctxMain); err != nil {
		cancelMain()
		return err
	}

	if err := a.initEventLog(ctxMain); err != nil {
		cancelMain()
		return err
	}

	defer func() {
		slog.Info("shutting down...")
		a.stopListener()
		a.close(ctx)
		cancelMain()
	}()

	ctxListen, cancelListen := WithSignal(ctxMain)
	defer cancelListen()

	a.initService()
	if err := a.listen(); err != nil {
		slog.Error("failed to listen",
			"error", err)
	}

	<-ctxListen.Done()

	return nil
}

func (a *Application) initEventLog(ctx context.Context) error {
	var tIDs []value.TenantID
	for t, err := range a.services.ts.List(ctx, 0, 0) {
		if err != nil {
			return err
		}

		tIDs = append(tIDs, t.ID())
	}

	return a.services.es.InitGlobalVersions(ctx, a.fsys, tIDs)
}

func (a *Application) listen() error {
	return listener.Listen(a.conf.Port, a.s)
}

func (a *Application) stopListener() {
	slog.Info("closing gRPC connections...")
	a.s.GracefulStop()
	slog.Info("closing gRPC connections complete")
}

// Closing Process
//   - close gRPC requests
//   - Cancel Context
//   - Flush Log
//   - Flush Buffers
func (a *Application) close(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	slog.Info("flushing...")
	if err := a.txManager.Shutdown(ctx, time.Now()); err != nil {
		slog.Error("unable to flush",
			"error", err)
	} else {
		slog.Info("flushing complete")
	}

	if err := a.ot.shutdown(ctx); err != nil {
		slog.Error("unable to close open telemtry",
			"error", err)
	}

	slog.Info("shutdown complete")
}

func (a *Application) initService() {
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
		a.services.es,
		a.services.os,
		a.services.ts,
	))

	a.s = s
}

type repositories struct {
}

func (a *Application) createRepositories() (*repositories, error) {
	return &repositories{}, nil
}
