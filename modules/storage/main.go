package main

import (
	"context"
	"fmt"

	"github.com/sjohnsonaz/dart-grpc-client/service/util/storage/disk"
	"github.com/sjohnsonaz/dart-grpc-client/service/util/storage/orchestrator"
)

func main() {
	err := writeDisk()
	if err != nil {
		fmt.Println(err)
	}
}

func writeDisk() error {
	d := disk.NewDisk()
	err := d.Open("mydisk.db", "myindex.db")
	if err != nil {
		return err
	}

	for _, value := range []string{
		"test0",
		"test1",
		"test2",
		"test3",
		"test4",
		"test5",
	} {
		index, err := d.WriteRecord(value)
		if err != nil {
			return err
		}

		fmt.Println(index)
	}

	// for r, _ := range d.Replay() {
	// 	fmt.Println("replay", r)
	// }

	// fmt.Println("finished replay")

	// for r, _ := range d.Replay() {
	// 	fmt.Println("replay", r)
	// }

	r, err := d.Get(3)
	if err != nil {
		return err
	}

	fmt.Println("get", r)

	r, err = d.Get(2)
	if err != nil {
		return err
	}

	fmt.Println("get", r)

	return nil
}

type writeAheadLog struct {
	watermark int
	records   []*record
}

type record struct{}

func (l *writeAheadLog) write(r *record) {
	l.records = append(l.records, r)
}

func (l *writeAheadLog) updateWatermark(w int) {
	l.watermark = w
}

func run() {
	o := orchestrator.NewOrchestrator()
	o.Run(context.Background())
}
