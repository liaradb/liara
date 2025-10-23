package main

import (
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

	err := application.Run()
	if err != nil {
		log.Fatal(err)
	}
}
