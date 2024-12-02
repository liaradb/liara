package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cardboardrobots/liara/esgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcConn, err := grpc.NewClient("localhost:50055",
		grpc.WithTransportCredentials(
			insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	es := esgrpc.NewEventSourceGRPC(grpcConn)
	count := 0
	for _, err := range es.ListTenants(context.Background()) {
		if err != nil {
			log.Fatal(err)
		}
		count++
	}
	if count == 0 {
		fmt.Println("no tenants")
	}
}
