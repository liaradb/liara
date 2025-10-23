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
	"google.golang.org/grpc"
)

func Run() error {
	conf, err := LoadConfig()
	if err != nil {
		return err
	}

	listener.Listen(context.Background(), conf.Port, conf.Port+1,
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
