package main

import (
	"fmt"
	"testing"
)

func TestWrite(t *testing.T) {
	err := writeDisk()
	if err != nil {
		fmt.Println(err)
	}
}

func TestWriteAheadLog(t *testing.T) {
	l := &writeAheadLog{}
	l.write(&record{})
	l.updateWatermark(0)
}

func TestOrchestrator(t *testing.T) {
	t.Skip()
	run()
}

func TestDisk(t *testing.T) {
	err := writeDisk()
	if err != nil {
		t.Error(err)
	}
}
