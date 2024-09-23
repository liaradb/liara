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
	service := listener.NewServerBuilder().
		AddUnary(
			listener.LogGRPC(false),
			listener.ErrorInterceptor(base.GetStatusCodeGRPC),
		).
		Build()

	pb.RegisterEventSourceServiceServer(service, &controller.EventSourceController{})

	return service
}
