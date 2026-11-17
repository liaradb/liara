package record

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type Action uint32

const ActionSize = 4

const (
	ActionShutdown     Action = 1
	ActionCheckpoint   Action = 2
	ActionStart        Action = 3
	ActionCommit       Action = 4
	ActionRollback     Action = 5
	ActionAppend       Action = 6
	ActionInsert       Action = 7
	ActionRemove       Action = 8
	ActionUpdate       Action = 9
	ActionClear        Action = 10
	ActionUpdateHeader Action = 11
)

const (
	actionNameShutdown     = "Shutdown"
	actionNameCheckpoint   = "Checkpoint"
	actionNameStart        = "Start"
	actionNameCommit       = "Commit"
	actionNameRollback     = "Rollback"
	actionNameAppend       = "Append"
	actionNameInsert       = "Insert"
	actionNameRemove       = "Remove"
	actionNameUpdate       = "Update"
	actionNameClear        = "Clear"
	actionNameUpdateHeader = "UpdateHeader"
	actionNameUnknown      = ""
)

func (Action) Size() int { return ActionSize }

func (a Action) Write(w io.Writer) error {
	return raw.WriteInt32(w, a)
}

func (a *Action) Read(r io.Reader) error {
	return raw.ReadInt32(r, a)
}

func (a Action) String() string {
	switch a {
	case ActionShutdown:
		return actionNameShutdown
	case ActionCheckpoint:
		return actionNameCheckpoint
	case ActionStart:
		return actionNameStart
	case ActionCommit:
		return actionNameCommit
	case ActionRollback:
		return actionNameRollback
	case ActionAppend:
		return actionNameAppend
	case ActionInsert:
		return actionNameInsert
	case ActionRemove:
		return actionNameRemove
	case ActionUpdate:
		return actionNameUpdate
	case ActionClear:
		return actionNameClear
	case ActionUpdateHeader:
		return actionNameUpdateHeader
	default:
		return actionNameUnknown
	}
}
