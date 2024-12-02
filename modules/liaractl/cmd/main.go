package main

import (
	"log"
	"os"

	"github.com/cardboardrobots/liaractl/handler"
)

func main() {
	ch, err := handler.NewCommandHandler("localhost:50055")
	if err != nil {
		log.Fatal(err)
	}

	if err := ch.Handle(os.Args); err != nil {
		log.Fatal(err)
	}
}
