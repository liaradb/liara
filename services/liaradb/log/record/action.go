package record

import (
	"encoding/binary"
	"io"
)

type Action uint32

const ActionSize = 4

const (
	ActionCheckpoint Action = 1
	ActionCommit     Action = 2
	ActionRollback   Action = 3
	ActionInsert     Action = 4
	ActionUpdate     Action = 5
	ActionRemove     Action = 6
)

func (a Action) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, a)
}

func (a *Action) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, a)
}

func (a Action) String() string {
	switch a {
	case ActionCheckpoint:
		return "Checkpoint"
	case ActionCommit:
		return "Commit"
	case ActionRollback:
		return "Rollback"
	case ActionInsert:
		return "Insert"
	case ActionRemove:
		return "Remove"
	case ActionUpdate:
		return "Update"
	default:
		return "Unknown"
	}
}
