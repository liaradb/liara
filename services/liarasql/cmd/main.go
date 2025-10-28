package main

import (
	"flag"
	"log"

	"github.com/cardboardrobots/liarasql/application"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

var (
	_ = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	log.SetPrefix("[liara]\t")
	log.Println("started...")

	err := application.Run()
	if err != nil {
		log.Fatal(err)
	}
}
