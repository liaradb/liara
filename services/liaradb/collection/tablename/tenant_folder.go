package tablename

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/filecache"
)

var eventLogRegexp = regexp.MustCompile("event_log--([0-9a-f]*).col")

type TenantFolder struct {
	tenantID value.TenantID
}

func NewTenantFolder(tid value.TenantID) TenantFolder {
	return TenantFolder{tid}
}

func NewTenantFolderFromString(s string) (TenantFolder, error) {
	tid, err := value.NewTenantIDFromString(s)
	if err != nil {
		return TenantFolder{}, err
	}

	return TenantFolder{tid}, nil
}

func (tf TenantFolder) TenantID() value.TenantID { return tf.tenantID }
func (tf TenantFolder) String() string           { return tf.tenantID.String() }

func (tf TenantFolder) GlobalTable(name string) GlobalTableName {
	return newGlobalTableName(tf, name)
}

func (tf TenantFolder) PartitionedTable(name string, pid value.PartitionID) PartitionedTableName {
	return newPartitionedTable(tf, name, pid)
}

func (tf TenantFolder) EventLog(pid value.PartitionID) PartitionedTableName {
	return tf.PartitionedTable("event_log", pid)
}

func (tf TenantFolder) RequestLog() GlobalTableName {
	return tf.GlobalTable("request_log")
}

func (tf TenantFolder) Outbox() GlobalTableName {
	return tf.GlobalTable("outbox")
}

func (tf TenantFolder) Tenant() GlobalTableName {
	return tf.GlobalTable("tenant")
}

func (tf TenantFolder) Partitions(fsys filecache.FileSystem, dir string) ([]value.PartitionID, error) {
	dirs, err := fsys.ReadDir(path.Join(dir, tf.String()))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = nil
		}
		return nil, err
	}

	pids := make([]value.PartitionID, 0, len(dirs))
	for _, dir := range dirs {
		num, ok := tf.partitionString(dir)
		if !ok {
			continue
		}

		i, err := strconv.Atoi(num)
		if err != nil {
			return nil, err
		}

		pids = append(pids, value.NewPartitionID(int32(i)))
	}

	return pids, nil
}

func (tf TenantFolder) partitionString(dir fs.DirEntry) (string, bool) {
	if dir.IsDir() {
		return "", false
	}

	n := dir.Name()
	if !IsCollection(n) {
		return "", false
	}

	_, fileName := path.Split(n)
	if !strings.HasPrefix(fileName, "event_log") {
		return "", false
	}

	matches := eventLogRegexp.FindStringSubmatch(fileName)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}
