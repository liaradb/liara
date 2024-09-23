package main

import (
	"context"
	"log"
	"net/http"

	"github.com/cardboardrobots/eventsource"
	pb "github.com/cardboardrobots/eventsource_go/generated"
	"github.com/cardboardrobots/liara_service/config"
	"github.com/cardboardrobots/liara_service/feature/base"
	"github.com/cardboardrobots/liara_service/feature/eventsource/controller"
	"github.com/cardboardrobots/liara_service/feature/eventsource/service"
	"github.com/cardboardrobots/listener"
	"google.golang.org/grpc"
)

func main() {
	log.SetPrefix("[liara]\t")
	log.Println("started...")

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	listener.Listen(context.Background(), conf.Port, conf.Port+1,
		http.NewServeMux(),
		initService())

	_ = eventsource.MockEventSource{}
}

func initService() *grpc.Server {
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

	pb.RegisterEventSourceServiceServer(s, controller.NewEventSourceController(
		service.NewEventService(&eventsource.MockEventSource{}),
	))

	return s
}
