package tenantfolder

import (
	"fmt"
	"io/fs"
	"path"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/util/testing/filetesting"
)

func TestTenantFolder(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	tf := NewTenantFolder(tid)

	if i := tf.TenantID(); i != tid {
		t.Errorf("incorrect tenant id: %v, expected: %v", i, tid)
	}

	if s := tf.String(); s != tid.String() {
		t.Errorf("incorrect string: %v, expected: %v", s, tid.String())
	}
}

func TestTenantFolder_NewTenantFolderFromString(t *testing.T) {
	t.Parallel()

	_, err := NewTenantFolderFromString("abcde")
	if err == nil {
		t.Error("should return error")
	}

	tid := value.NewTenantID()
	tidString := tid.String()
	tf, err := NewTenantFolderFromString(tidString)
	if err != nil {
		t.Error(err)
	}

	if i := tf.TenantID(); i != tid {
		t.Errorf("incorrect tenant id: %v, expected: %v", i, tid)
	}

	if s := tf.String(); s != tidString {
		t.Errorf("incorrect string: %v, expected: %v", s, tidString)
	}
}

func TestTenantFolder_EventLog(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	tf := NewTenantFolder(tid)
	pid := value.NewPartitionID(1)

	el := tf.EventLog(pid)

	want := fmt.Sprintf("%v/event_log--%08x.col", tid, 1)
	if s := el.Collection().String(); s != want {
		t.Errorf("incorrect string: %v, expected: %v", s, want)
	}
}

func TestTenantFolder_Outbox(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	tf := NewTenantFolder(tid)

	ob := tf.Outbox()

	want := fmt.Sprintf("%v/outbox.col", tid)
	if s := ob.Collection().String(); s != want {
		t.Errorf("incorrect string: %v, expected: %v", s, want)
	}
}

func TestTenantFolder_RequestLog(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	tf := NewTenantFolder(tid)

	rl := tf.RequestLog()

	want := fmt.Sprintf("%v/request_log.col", tid)
	if s := rl.Collection().String(); s != want {
		t.Errorf("incorrect string: %v, expected: %v", s, want)
	}
}

func TestTenantFolder_Tenant(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	tf := NewTenantFolder(tid)

	tn := tf.Tenant()

	want := fmt.Sprintf("%v/tenant.col", tid)
	if s := tn.Collection().String(); s != want {
		t.Errorf("incorrect string: %v, expected: %v", s, want)
	}
}

func TestTenantFolder_Partitions(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	tf := NewTenantFolder(tid)

	fsys := filetesting.New(fstest.MapFS{
		path.Join("dir", tf.EventLog(value.NewPartitionID(1)).Collection().String()): &fstest.MapFile{},
		path.Join("dir", tf.EventLog(value.NewPartitionID(2)).Collection().String()): &fstest.MapFile{},
		path.Join("dir", tf.EventLog(value.NewPartitionID(2)).Index(0).String()):     &fstest.MapFile{},
		path.Join("dir", tf.Outbox().Collection().String()):                          &fstest.MapFile{},
		path.Join("dir", tf.String(), "event_log.col"):                               &fstest.MapFile{},
		path.Join("dir", tf.String(), "nested"):                                      &fstest.MapFile{Mode: fs.ModeDir},
	})

	pids, err := tf.Partitions(fsys, "dir")
	if err != nil {
		t.Fatal(err)
	}

	want := []value.PartitionID{value.NewPartitionID(1), value.NewPartitionID(2)}
	if !slices.Equal(pids, want) {
		t.Errorf("incorrect partition ids: %v, expected: %v", pids, want)
	}
}

func TestTenantFolder_Partitions__Empty(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	tf := NewTenantFolder(tid)

	fsys := filetesting.New(fstest.MapFS{
		path.Join("dir"): &fstest.MapFile{Mode: fs.ModeDir},
	})

	pids, err := tf.Partitions(fsys, "dir")
	if err != nil {
		t.Fatal(err)
	}

	want := []value.PartitionID{}
	if !slices.Equal(pids, want) {
		t.Errorf("incorrect partition ids: %v, expected: %v", pids, want)
	}
}

func TestTenantFolder_Partitions__Error(t *testing.T) {
	t.Parallel()

	tid := value.NewTenantID()
	tf := NewTenantFolder(tid)

	fsys := filetesting.New(fstest.MapFS{
		path.Join("dir", tf.String(), "event_log--a.col"): &fstest.MapFile{},
	})

	_, err := tf.Partitions(fsys, "dir")
	if err == nil {
		t.Fatal("should return error")
	}
}
