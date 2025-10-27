package main

import (
	"context"
	"flag"
	"log/slog"

	"github.com/liaradb/liaradb/application"
)

var (
	_ = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	// log.SetPrefix("[liaradb]\t")
	slog.Info("started...")

	conf, err := application.LoadConfig()
	if err != nil {
		slog.Error("main", "error", err)
		return
	}

	a := application.New(conf)
	err = a.Run(context.Background())
	if err != nil {
		slog.Error("main", "error", err)
	}
}
