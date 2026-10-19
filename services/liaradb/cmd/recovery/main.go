package main

import (
	"flag"
	"fmt"
	"log/slog"
	"path"

	"github.com/liaradb/liaradb/application"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery/segment"
	"github.com/liaradb/liaradb/transaction/log"
)

var (
	_ = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		slog.Error("recovery", "error", err)
	}
}

func run() error {
	conf, err := application.LoadConfig()
	if err != nil {
		return err
	}

	segmentSize := 1024
	fsys := filecache.New()

	log := log.New(
		int64(conf.BlockSize),
		segment.PageID(segmentSize),
		int64(conf.RecordSize),
		conf.WriteQueueSize,
		fsys,
		path.Join(conf.Directory, "log"))

	// TODO: Change back to log.Recover()
	// it, err := log.Recover()
	it, err := log.Iterate()
	if err != nil {
		slog.Error("recovery", "error", err)
	}

	for rc := range it {
		if err != nil {
			slog.Error("recovery", "error", err)
		}

		fmt.Printf("%v\t%v\t%-12s\t%v\n",
			rc.LogSequenceNumber(),
			rc.TransactionID(),
			fmt.Sprintf("<%v>", rc.Action()),
			rc.RecordLocator())
	}

	return nil
}
