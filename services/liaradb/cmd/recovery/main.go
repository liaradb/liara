package main

import (
	"flag"
	"fmt"
	"log/slog"
	"path"

	"github.com/liaradb/liaradb/application"
	"github.com/liaradb/liaradb/file/disk"
	"github.com/liaradb/liaradb/recovery"
	"github.com/liaradb/liaradb/recovery/action"
	"github.com/liaradb/liaradb/recovery/record"
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
	fsys := &disk.FileSystem{}

	log := recovery.NewLog(
		int64(conf.BlockSize),
		action.PageID(segmentSize),
		fsys,
		path.Join(conf.Directory, "log"))

	for rc, err := range log.Iterate(record.NewLogSequenceNumber(0)) {
		if err != nil {
			slog.Error("recovery", "error", err)
		}

		fmt.Printf("%v\t%v\t<%v>\n", rc.LogSequenceNumber(), rc.TransactionID(), rc.Action())
	}

	return nil
}
