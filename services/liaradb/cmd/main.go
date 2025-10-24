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

	conf, err := application.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	a := application.New(conf)
	err = a.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
