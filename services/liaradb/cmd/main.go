package main

import (
	"context"
	"flag"
	"log"

	"github.com/liaradb/liaradb/application"
)

var (
	_ = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	log.SetPrefix("[liaradb]\t")
	log.Println("started...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a := application.New(200, 4096)
	err := a.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
