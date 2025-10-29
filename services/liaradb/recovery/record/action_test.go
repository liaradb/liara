package record

import (
	"io"
	"testing"

	"github.com/cardboardrobots/assert"
)

func TestAction(t *testing.T) {
	t.Parallel()

	r, w := assert.NewReaderWriter()

	var a Action = ActionCheckpoint
	if err := a.Write(w); err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	var a2 Action
	if err := a2.Read(r); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if a != a2 {
		t.Errorf("incorrect value: %v, expected: %v", a2, a)
	}
}

func TestAction_String(t *testing.T) {
	t.Parallel()

	for message, c := range map[string]struct {
		action Action
		result string
	}{
		"should handle zero": {
			action: 0,
			result: "Unknown"},
		"should handle checkpoint": {
			action: ActionCheckpoint,
			result: "Checkpoint"},
		"should handle commit": {
			action: ActionCommit,
			result: "Commit"},
		"should handle rollback": {
			action: ActionRollback,
			result: "Rollback"},
		"should handle insert": {
			action: ActionInsert,
			result: "Insert"},
		"should handle remove": {
			action: ActionRemove,
			result: "Remove"},
		"should handle update": {
			action: ActionUpdate,
			result: "Update"},
	} {
		t.Run(message, func(t *testing.T) {
			t.Parallel()
			if r := c.action.String(); r != c.result {
				t.Errorf("incorrect string: %v, expected: %v", r, c.result)
			}
		})
	}
}
