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

	if err := run(); err != nil {
		slog.Error("main", "error", err)
	}
}

func run() error {
	conf, err := application.LoadConfig()
	if err != nil {
		return err
	}

	a := application.New(conf)
	return a.Run(context.Background())
}
