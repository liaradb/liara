package record

import (
	"io"

	"github.com/liaradb/liaradb/encoder/raw"
)

type Action uint32

const ActionSize = 4

const (
	ActionCheckpoint   Action = 1
	ActionStart        Action = 2
	ActionCommit       Action = 3
	ActionRollback     Action = 4
	ActionAppend       Action = 5
	ActionInsert       Action = 6
	ActionRemove       Action = 7
	ActionUpdate       Action = 8
	ActionClear        Action = 9
	ActionUpdateHeader Action = 10
)

const (
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
