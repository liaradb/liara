package tablename

import (
	"testing"

	"github.com/liaradb/liaradb/domain/value"
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
