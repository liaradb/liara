package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cardboardrobots/liara/esgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type object string
type command string

const (
	objectTenant object  = "tenant"
	objectEvent  object  = "event"
	commandList  command = "list"
)

func main() {
	cmd := getArgs()
	switch cmd {
	case objectTenant:
		listTenants()
	case objectEvent:
		listEvents()
	default:
		fmt.Println("no command")
		break
	}
}

func listTenants() {
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

func listEvents() {
	grpcConn, err := grpc.NewClient("localhost:50055",
		grpc.WithTransportCredentials(
			insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	es := esgrpc.NewEventSourceGRPC(grpcConn)
	count := 0
	for event, err := range es.GetAfterGlobalVersion(context.Background(), 0, nil, 0) {
		if err != nil {
			log.Fatal(err)
		}

		var data = make(map[string]any)
		_ = json.Unmarshal(event.Data, &data)

		log.Printf("%v\n", data)
		count++
	}
	if count == 0 {
		fmt.Println("no events")
	}
}

func getArgs() object {
	args := os.Args
	switch len(args) {
	case 0:
		fallthrough
	case 1:
		return ""
	default:
		return object(args[1])
	}
}
