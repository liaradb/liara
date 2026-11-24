package record

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type Action uint32

const ActionSize = 4

const (
	ActionShutdown     Action = 1
	ActionRecover      Action = 2
	ActionCheckpoint   Action = 3
	ActionStart        Action = 4
	ActionCommit       Action = 5
	ActionRollback     Action = 6
	ActionAppend       Action = 7
	ActionInsert       Action = 8
	ActionRemove       Action = 9
	ActionUpdate       Action = 10
	ActionClear        Action = 11
	ActionUpdateHeader Action = 12
)

const (
	actionNameShutdown     = "Shutdown"
	actionNameRecover      = "Recover"
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
	case ActionRecover:
		return actionNameRecover
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
