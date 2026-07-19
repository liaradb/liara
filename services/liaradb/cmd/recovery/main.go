package main

import (
	"flag"
	"fmt"
	"log/slog"
	"path"

	"github.com/liaradb/liaradb/application"
	"github.com/liaradb/liaradb/filecache"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/segment"
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

	log := recovery.NewLog(
		int64(conf.BlockSize),
		segment.PageID(segmentSize),
		int64(conf.RecordSize),
		conf.WriteQueueSize,
		fsys,
		path.Join(conf.Directory, "log"))

	it, err := log.Recover()
	if err != nil {
		slog.Error("recovery", "error", err)
	}

	for rc := range it {
		if err != nil {
			slog.Error("recovery", "error", err)
		}

		fmt.Printf("%v\t%v\t<%v>\n", rc.LogSequenceNumber(), rc.TransactionID(), rc.Action())
	}

	return nil
}
